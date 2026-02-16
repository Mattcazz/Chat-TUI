package user

import (
	"encoding/json"
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
		utils.WriteJSONError(w, http.StatusBadRequest, err)
		return
	}

	if req.Username == "" || req.PublicKey == "" {
		utils.WriteJsonMsg(w, http.StatusBadRequest, "Username and public key are required")
		return
	}

	user, err := h.service.CreateUser(ctx, req.PublicKey, req.Username)
	if err != nil {
		utils.WriteJSONError(w, http.StatusBadRequest, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, user)
}

func (h *Handler) login(w http.ResponseWriter, r *http.Request) {
	var req pkg.LoginRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteJSONError(w, http.StatusBadRequest, err)
		return
	}

	if req.PublicKey == "" {
		utils.WriteJsonMsg(w, http.StatusBadRequest, "Public key is required")
		return
	}

	if req.Signature == "" {
		nonce, err := h.service.GenerateChallenge(r.Context(), req.PublicKey)
		if err != nil {
			if IsUserDoesNotExistError(err) {
				utils.WriteJsonMsg(w, http.StatusTemporaryRedirect, err.Error())
				return
			}
			utils.WriteJSONError(w, http.StatusBadRequest, err)
			return
		}

		resp := pkg.ChallengeResponse{Nonce: nonce}

		utils.WriteJSON(w, http.StatusOK, resp)
		return
	}

	token, err := h.service.VerifyAndLogin(r.Context(), req.PublicKey, req.Signature)
	if err != nil {
		utils.WriteJSONError(w, http.StatusBadRequest, err)
		return
	}

	resp := pkg.LoginResponse{Token: token}

	utils.WriteJSON(w, http.StatusOK, resp)
}

func (h *Handler) getInbox(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(utils.CtxKeyUserID)

	inbox, err := h.service.GetInbox(r.Context(), userID.(int64))
	if err != nil {
		utils.WriteJSONError(w, http.StatusBadRequest, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, inbox)
}

func (h *Handler) deleteUser(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(utils.CtxKeyUserID)

	err := h.service.DeleteUser(r.Context(), userID.(int64))
	if err != nil {
		utils.WriteJSONError(w, http.StatusBadRequest, err)
		return
	}

	utils.WriteJsonMsg(w, http.StatusOK, "User deleted")
}

func (h *Handler) getContacts(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(utils.CtxKeyUserID)

	contacts, err := h.service.GetContacts(r.Context(), userID.(int64))
	if err != nil {
		utils.WriteJSONError(w, http.StatusBadRequest, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, contacts)
}

func (h *Handler) postContactRequest(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(utils.CtxKeyUserID)

	var req pkg.PostContactRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteJSONError(w, http.StatusBadRequest, err)
		return
	}

	err := h.service.ContactRequest(r.Context(), userID.(int64), req.PublicKey, req.Nickname)
	if err != nil {
		utils.WriteJSONError(w, http.StatusBadRequest, err)
		return
	}

	utils.WriteJsonMsg(w, http.StatusOK, "Contact request sent")
}

func (h *Handler) getContactRequests(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(utils.CtxKeyUserID)

	contacts, err := h.service.GetContactRequests(r.Context(), userID.(int64))
	if err != nil {
		utils.WriteJSONError(w, http.StatusBadRequest, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, contacts)
}

func (h *Handler) blockContact(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(utils.CtxKeyUserID)

	contactIDStr := chi.URLParam(r, "contact_id")

	contactID, err := strconv.Atoi(contactIDStr)
	if err != nil {
		utils.WriteJSONError(w, http.StatusBadRequest, err)
		return
	}

	err = h.service.BlockContact(r.Context(), userID.(int64), int64(contactID))
	if err != nil {
		utils.WriteJSONError(w, http.StatusBadRequest, err)
		return
	}

	utils.WriteJsonMsg(w, http.StatusOK, "Contact blocked")
}

func (h *Handler) unblockContact(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(utils.CtxKeyUserID)

	contactIDStr := chi.URLParam(r, "contact_id")

	contactID, err := strconv.Atoi(contactIDStr)
	if err != nil {
		utils.WriteJSONError(w, http.StatusBadRequest, err)
		return
	}

	err = h.service.UnblockContact(r.Context(), userID.(int64), int64(contactID))
	if err != nil {
		utils.WriteJSONError(w, http.StatusBadRequest, err)
		return
	}

	utils.WriteJsonMsg(w, http.StatusOK, "Contact unblocked")
}
