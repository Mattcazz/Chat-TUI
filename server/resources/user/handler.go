package user

import (
	"encoding/json"
	"fmt"
	"net/http"

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
	r.Post("/", h.userChallenge)
	r.Post("/login", h.login)
	r.Post("/register", h.registerUser)
	r.Route("/contacts", func(r chi.Router) {
		r.Get("/", middleware.JWTAuth(h.getContacts))
		r.Post("/request", middleware.JWTAuth(h.postContactRequest))
		r.Patch("/{contact_id}", middleware.JWTAuth(h.patchContact))
	})
}

func (h *Handler) registerUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req pkg.RegisterRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteJSONError(w, http.StatusBadRequest, err)
		return
	}

	user, err := h.service.CreateUser(ctx, req.PublicKey, req.Username)

	if err != nil {
		utils.WriteJSONError(w, http.StatusBadRequest, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, user)
}

func (h *Handler) userChallenge(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req pkg.ChallengeRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteJSONError(w, http.StatusBadRequest, err)
		return
	}

	nonce, err := h.service.GenerateChallenge(ctx, req.PublicKey)

	if err != nil {
		if IsUserDoesNotExistError(err) {
			utils.WriteJSONError(w, http.StatusPermanentRedirect, err)
			return
		}
		utils.WriteJSONError(w, http.StatusBadRequest, err)
		return
	}

	resp := pkg.ChallengeResponse{Nonce: nonce}

	utils.WriteJSON(w, http.StatusOK, resp)
}

func (h *Handler) login(w http.ResponseWriter, r *http.Request) {
	var req pkg.LoginRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteJSONError(w, http.StatusBadRequest, err)
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

}

func (h *Handler) getContacts(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(utils.CtxKeyUserID)

	if userID == nil {
		utils.WriteJSONError(w, http.StatusBadRequest, fmt.Errorf("User ID not found in context"))
		return
	}

	utils.WriteJSON(w, http.StatusOK, userID)
}

func (h *Handler) postContactRequest(w http.ResponseWriter, r *http.Request) {

}

func (h *Handler) patchContact(w http.ResponseWriter, r *http.Request) {
}
