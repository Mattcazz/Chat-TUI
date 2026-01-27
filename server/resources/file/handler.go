package file

import "github.com/go-chi/chi/v5"

type Handler struct {
	Service *Service
}

func NewHandler(s *Service) *Handler {
	return &Handler{
		Service: s,
	}
}

func (h *Handler) RegisterRoutes(r *chi.Mux) {

}
