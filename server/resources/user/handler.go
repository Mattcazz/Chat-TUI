package user

import (
	"net/http"

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

	r.Route("/contacts", func(r chi.Router) {
		r.Get("/", h.getContacts)
		r.Post("/request", h.postContactRequest)
		r.Patch("/{user_id}", h.patchContact)
	})
}

func (h *Handler) getInbox(w http.ResponseWriter, r *http.Request) {

}

func (h *Handler) getContacts(w http.ResponseWriter, r *http.Request) {

}

func (h *Handler) postContactRequest(w http.ResponseWriter, r *http.Request) {

}

func (h *Handler) patchContact(w http.ResponseWriter, r *http.Request) {
}
