package file

import "github.com/go-chi/chi/v5"

type Handler struct {
	file_store FileStore
}

func NewHandler(fs FileStore) *Handler {
	return &Handler{
		file_store: fs,
	}
}

func (h *Handler) RegisterRoutes(r *chi.Mux) {

}
