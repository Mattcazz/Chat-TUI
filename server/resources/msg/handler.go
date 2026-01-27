package msg

import "github.com/go-chi/chi/v5"

type Handler struct {
	conv_service *Service
}

func NewHandler(s *Service) *Handler {
	return &Handler{
		conv_service: s,
	}
}

func (h *Handler) RegisterRoutes(r *chi.Mux) {

}
