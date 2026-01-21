package msg

import "net/http"

type Handler struct {
	msg_store MsgStore
}

func NewHandler(ms MsgStore) *Handler {
	return &Handler{
		msg_store: ms,
	}
}

func (h *Handler) RegisterRoutes(m *http.ServeMux) {

}
