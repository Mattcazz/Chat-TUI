package contact

import "github.com/go-chi/chi/v5"

type Handler struct {
	contact_store ContactStore
}

func NewHandler(cs ContactStore) *Handler {
	return &Handler{
		contact_store: cs,
	}
}

func (h *Handler) RegisterRoutes(r *chi.Mux) {

}
