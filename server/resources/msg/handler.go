package msg

import "github.com/go-chi/chi/v5"

type Handler struct {
	convService *Service
}

func NewHandler(s *Service) *Handler {
	return &Handler{
		convService: s,
	}
}

func (h *Handler) RegisterRoutes(r *chi.Mux) {

}
