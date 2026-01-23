package user

import "github.com/go-chi/chi/v5"

type Handler struct {
	user_store UserStore
}

func NewHandler(cs UserStore) *Handler {
	return &Handler{
		user_store: cs,
	}
}

func (h *Handler) RegisterRoutes(r *chi.Mux) {

}
