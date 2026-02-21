package chat

import (
	"encoding/json"
	"net/http"

	"github.com/Mattcazz/Chat-TUI/pkg"
	"github.com/Mattcazz/Chat-TUI/server/resources/middleware"
	"github.com/Mattcazz/Chat-TUI/server/utils"
	"github.com/docker/docker/volume/service"
	"github.com/go-chi/chi/v5"
)

type Handler struct {
	convService *Service
}

func NewHandler(s *Service) *Handler {
	return &Handler{
		convService: s,
	}
}

func (h *Handler) RegisterRoutes(r *chi.Mux) {
	r.Route("/chat", func(r chi.Router) {
		r.Post("/", middleware.JWTAuth(h.postConversation))
		r.Get("/{conversation_id}", middleware.JWTAuth(h.getConversation))
		r.Delete("/{conversation_id}", middleware.JWTAuth(h.deleteConversation))
		r.Post("/{conversation_id}/message", middleware.JWTAuth(h.postMessageInConversation))
	})
}

func (h *Handler) getConversation(w http.ResponseWriter, r *http.Request) {
}

func (h *Handler) postConversation(w http.ResponseWriter, r *http.Request) {
}

func (h *Handler) deleteConversation(w http.ResponseWriter, r *http.Request) {
}

func (h *Handler) postMessageInConversation(w http.ResponseWriter, r *http.Request) {
	senderId := r.Context().Value(utils.CtxKeyUserID)

	var msgReq pkg.SendMsgRequest

	if err := json.NewEncoder(w).Encode(msgReq); err != nil {
		utils.WriteJSONError(w, http.StatusBadRequest, err)
	}

	conversation_id := r.chi.URLParam(r, "conversation_id")

	if err = h.convService.postConversationMsg(ctx context.Context, sender_id int64, conv_id int64, content string)
}
