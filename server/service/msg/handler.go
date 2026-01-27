package msg

import "github.com/go-chi/chi/v5"

type Handler struct {
	msg_store MsgStore
}

func NewHandler(ms MsgStore) *Handler {
	return &Handler{
		msg_store: ms,
	}
}

func (h *Handler) RegisterRoutes(r *chi.Mux) {
}
