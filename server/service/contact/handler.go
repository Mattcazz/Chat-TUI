package contact

import "net/http"

type Handler struct {
	contact_store ContactStore
}

func NewHandler(cs ContactStore) *Handler {
	return &Handler{
		contact_store: cs,
	}
}

func (h *Handler) RegisterRoutes(m *http.ServeMux) {

}
