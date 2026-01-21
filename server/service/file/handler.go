package file

import "net/http"

type Handler struct {
	file_store FileStore
}

func NewHandler(fs FileStore) *Handler {
	return &Handler{
		file_store: fs,
	}
}

func (h *Handler) RegisterRoutes(m *http.ServeMux) {

}
