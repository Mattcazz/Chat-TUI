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
	r.Route("/files", func(r chi.Router) {
		r.Post("/init", h.fileInit)
		r.Post("/upload/{session_id}/chunks", h.uploadChunk)
		r.Post("/upload/{session_id}/assemble", h.assembleFile)
		r.Get("/upload/{session_id}/status", h.statusCheck)
		r.Get("/download/{file_id}", h.downloadFile)
		r.Delete("/upload/{session_id}", h.cancelUpload)
	})
}

func (h *Handler) fileInit(w http.ResponseWriter, r *http.Request) {
}

func (h *Handler) uploadChunk(w http.ResponseWriter, r *http.Request) {
}

func (h *Handler) assembleFile(w http.ResponseWriter, r *http.Request) {
}

func (h *Handler) statusCheck(w http.ResponseWriter, r *http.Request) {
}

func (h *Handler) downloadFile(w http.ResponseWriter, r *http.Request) {
}

func (h *Handler) cancelUpload(w http.ResponseWriter, r *http.Request) {
}
