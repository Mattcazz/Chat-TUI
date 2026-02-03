package user

import (
	"encoding/json"
	"net/http"

	"github.com/Mattcazz/Chat-TUI/pkg"

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

	r.Get("/inbox", h.getInbox)
	r.Post("/", h.userChallenge)
	r.Post("/login", h.login)

	r.Route("/contacts", func(r chi.Router) {
		r.Get("/", h.getContacts)
		r.Post("/request", h.postContactRequest)
		r.Patch("/{contact_id}", h.patchContact)
	})
}

func (h *Handler) userChallenge(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req pkg.ChallengeRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteJSONError(w, 400, err)
		return
	}

	nonce, err := h.service.GenerateChallenge(ctx, req.PublicKey)

	if err != nil {
		utils.WriteJSONError(w, 404, err)
		return
	}

	resp := pkg.ChallengeResponse{Nonce: nonce}

	utils.WriteJSON(w, 200, resp)
}

func (h *Handler) login(w http.ResponseWriter, r *http.Request) {
	var req *pkg.LoginRequest

	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		utils.WriteJSONError(w, 400, err)
		return
	}

	token, err := h.service.VerifyAndLogin(r.Context(), req.PublicKey, req.Signature)

	if err != nil {
		utils.WriteJSONError(w, 400, err)
		return
	}

	resp := pkg.LoginResponse{Token: token}

	json.NewEncoder(w).Encode(resp)
}

func (h *Handler) getInbox(w http.ResponseWriter, r *http.Request) {

}

func (h *Handler) getContacts(w http.ResponseWriter, r *http.Request) {

}

func (h *Handler) postContactRequest(w http.ResponseWriter, r *http.Request) {

}

func (h *Handler) patchContact(w http.ResponseWriter, r *http.Request) {
}
