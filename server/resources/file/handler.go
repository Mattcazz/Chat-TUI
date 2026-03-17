package file

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Handler struct {
	Service *Service
}

func NewHandler(s *Service) *Handler {
	return &Handler{
		Service: s,
	}
}

func (h *Handler) RegisterRoutes(r *chi.Mux) {
	r.Route("/file", func(r chi.Router) {
		r.Post("/", h.uploadFile)
		r.Get("/{file_id}", h.downloadFile)
	})
}

func (h *Handler) uploadFile(w http.ResponseWriter, r *http.Request) {
}

func (h *Handler) downloadFile(w http.ResponseWriter, r *http.Request) {
}
