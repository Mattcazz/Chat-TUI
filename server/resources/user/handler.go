package user

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/Mattcazz/Chat-TUI/pkg"

	"github.com/Mattcazz/Chat-TUI/server/resources/middleware"
	"github.com/Mattcazz/Chat-TUI/server/utils"
	"github.com/go-chi/chi/v5"
)

type Handler struct {
	service *Service
}

func NewHandler(s *Service) *Handler {
	return &Handler{
		service: s,
	}
}

func (h *Handler) RegisterRoutes(r *chi.Mux) {
	r.Get("/inbox", middleware.JWTAuth(h.getInbox))
	r.Post("/login", h.login)
	r.Post("/register", h.registerUser)
	r.Delete("/delete", middleware.JWTAuth(h.deleteUser))
	r.Route("/contacts", func(r chi.Router) {
		r.Get("/", middleware.JWTAuth(h.getContacts))
		r.Post("/requests", middleware.JWTAuth(h.postContactRequest))
		r.Get("/requests", middleware.JWTAuth(h.getContactRequests))
		r.Post("/{contact_id}/block", middleware.JWTAuth(h.blockContact))
		r.Post("/{contact_id}/unblock", middleware.JWTAuth(h.unblockContact))
	})
}

func (h *Handler) registerUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req pkg.RegisterRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Failed to decode register request: %v", err)
		utils.WriteJSONError(w, http.StatusBadRequest, err)
		return
	}

	if req.Username == "" || req.PublicKey == "" {
		log.Println("Register request missing username or public key")
		utils.WriteJsonMsg(w, http.StatusBadRequest, "Username and public key are required")
		return
	}

	log.Printf("Attempting to register user with username: %s", req.Username)
	user, err := h.service.CreateUser(ctx, req.PublicKey, req.Username)
	if err != nil {
		log.Printf("Failed to create user with username %s: %v", req.Username, err)
		utils.WriteJSONError(w, http.StatusBadRequest, err)
		return
	}

	log.Printf("Successfully registered user with ID: %d, username: %s", user.ID, user.Username)
	utils.WriteJSON(w, http.StatusOK, user)
}

func (h *Handler) login(w http.ResponseWriter, r *http.Request) {
	var req pkg.LoginRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Failed to decode login request: %v", err)
		utils.WriteJSONError(w, http.StatusBadRequest, err)
		return
	}

	if req.PublicKey == "" {
		log.Println("Login request missing public key")
		utils.WriteJsonMsg(w, http.StatusBadRequest, "Public key is required")
		return
	}

	if req.Signature == "" {
		log.Println("Generating challenge for login (no signature provided)")
		nonce, err := h.service.GenerateChallenge(r.Context(), req.PublicKey)
		if err != nil {
			if IsUserDoesNotExistError(err) {
				log.Println("User does not exist for login challenge generation")
				utils.WriteJsonMsg(w, http.StatusTemporaryRedirect, err.Error())
				return
			}
			log.Printf("Failed to generate challenge: %v", err)
			utils.WriteJSONError(w, http.StatusBadRequest, err)
			return
		}

		log.Println("Challenge generated successfully")
		resp := pkg.ChallengeResponse{Nonce: nonce}

		utils.WriteJSON(w, http.StatusAccepted, resp)
		return
	}

	log.Println("Verifying login with signature")
	token, err := h.service.VerifyAndLogin(r.Context(), req.PublicKey, req.Signature)
	if err != nil {
		log.Printf("Failed to verify and login: %v", err)
		utils.WriteJSONError(w, http.StatusBadRequest, err)
		return
	}

	log.Println("User login verified successfully")
	resp := pkg.LoginResponse{Token: token}

	utils.WriteJSON(w, http.StatusOK, resp)
}

func (h *Handler) getInbox(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(utils.CtxKeyUserID)

	log.Printf("Fetching inbox for user ID: %d", userID.(int64))
	inbox, err := h.service.GetInbox(r.Context(), userID.(int64))
	if err != nil {
		log.Printf("Failed to retrieve inbox for user ID %d: %v", userID.(int64), err)
		utils.WriteJSONError(w, http.StatusBadRequest, err)
		return
	}

	log.Printf("Successfully retrieved inbox for user ID: %d", userID.(int64))
	utils.WriteJSON(w, http.StatusOK, inbox)
}

func (h *Handler) deleteUser(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(utils.CtxKeyUserID)

	log.Printf("Deleting user with ID: %d", userID.(int64))
	err := h.service.DeleteUser(r.Context(), userID.(int64))
	if err != nil {
		log.Printf("Failed to delete user with ID %d: %v", userID.(int64), err)
		utils.WriteJSONError(w, http.StatusBadRequest, err)
		return
	}

	log.Printf("Successfully deleted user with ID: %d", userID.(int64))
	utils.WriteJsonMsg(w, http.StatusOK, "User deleted")
}

func (h *Handler) getContacts(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(utils.CtxKeyUserID)

	log.Printf("Fetching contacts for user ID: %d", userID.(int64))
	contacts, err := h.service.GetContacts(r.Context(), userID.(int64))
	if err != nil {
		log.Printf("Failed to retrieve contacts for user ID %d: %v", userID.(int64), err)
		utils.WriteJSONError(w, http.StatusBadRequest, err)
		return
	}

	log.Printf("Successfully retrieved contacts for user ID: %d", userID.(int64))
	utils.WriteJSON(w, http.StatusOK, contacts)
}

func (h *Handler) postContactRequest(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(utils.CtxKeyUserID)

	var req pkg.PostContactRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Failed to decode contact request for user ID %d: %v", userID.(int64), err)
		utils.WriteJSONError(w, http.StatusBadRequest, err)
		return
	}

	log.Printf("Processing contact request from user ID %d with nickname: %s", userID.(int64), req.Nickname)
	err := h.service.ContactRequest(r.Context(), userID.(int64), req.PublicKey, req.Nickname)
	if err != nil {
		log.Printf("Failed to send contact request from user ID %d: %v", userID.(int64), err)
		utils.WriteJSONError(w, http.StatusBadRequest, err)
		return
	}

	log.Printf("Successfully sent contact request from user ID %d", userID.(int64))
	utils.WriteJsonMsg(w, http.StatusOK, "Contact request sent")
}

func (h *Handler) getContactRequests(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(utils.CtxKeyUserID)

	log.Printf("Fetching contact requests for user ID: %d", userID.(int64))
	contacts, err := h.service.GetContactRequests(r.Context(), userID.(int64))
	if err != nil {
		log.Printf("Failed to retrieve contact requests for user ID %d: %v", userID.(int64), err)
		utils.WriteJSONError(w, http.StatusBadRequest, err)
		return
	}

	log.Printf("Successfully retrieved contact requests for user ID: %d", userID.(int64))
	utils.WriteJSON(w, http.StatusOK, contacts)
}

func (h *Handler) blockContact(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(utils.CtxKeyUserID)

	contactIDStr := chi.URLParam(r, "contact_id")
	log.Printf("Blocking contact for user ID: %d, contact_id: %s", userID.(int64), contactIDStr)

	contactID, err := strconv.Atoi(contactIDStr)
	if err != nil {
		log.Printf("Failed to parse contact ID for user ID %d: %v", userID.(int64), err)
		utils.WriteJSONError(w, http.StatusBadRequest, err)
		return
	}

	err = h.service.BlockContact(r.Context(), userID.(int64), int64(contactID))
	if err != nil {
		log.Printf("Failed to block contact ID %d for user ID %d: %v", contactID, userID.(int64), err)
		utils.WriteJSONError(w, http.StatusBadRequest, err)
		return
	}

	log.Printf("Successfully blocked contact ID %d for user ID %d", contactID, userID.(int64))
	utils.WriteJsonMsg(w, http.StatusOK, "Contact blocked")
}

func (h *Handler) unblockContact(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(utils.CtxKeyUserID)

	contactIDStr := chi.URLParam(r, "contact_id")
	log.Printf("Unblocking contact for user ID: %d, contact_id: %s", userID.(int64), contactIDStr)

	contactID, err := strconv.Atoi(contactIDStr)
	if err != nil {
		log.Printf("Failed to parse contact ID for user ID %d: %v", userID.(int64), err)
		utils.WriteJSONError(w, http.StatusBadRequest, err)
		return
	}

	err = h.service.UnblockContact(r.Context(), userID.(int64), int64(contactID))
	if err != nil {
		log.Printf("Failed to unblock contact ID %d for user ID %d: %v", contactID, userID.(int64), err)
		utils.WriteJSONError(w, http.StatusBadRequest, err)
		return
	}

	log.Printf("Successfully unblocked contact ID %d for user ID %d", contactID, userID.(int64))
	utils.WriteJsonMsg(w, http.StatusOK, "Contact unblocked")
}
