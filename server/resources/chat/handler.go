package chat

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Mattcazz/Chat-TUI/pkg"
	"github.com/Mattcazz/Chat-TUI/server/resources/middleware"
	"github.com/Mattcazz/Chat-TUI/server/utils"
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
	r.Route("/conversation", func(r chi.Router) {
		r.Post("/", middleware.JWTAuth(h.postConversationDM))
		r.Get("/{conversation_id}", middleware.JWTAuth(h.getConversation))
		r.Delete("/{conversation_id}", middleware.JWTAuth(h.deleteConversation))
		r.Post("/{conversation_id}/message", middleware.JWTAuth(h.postMessageInConversation))
	})
}

func (h *Handler) getConversation(w http.ResponseWriter, r *http.Request) {
	conversationIDStr := chi.URLParam(r, "conversation_id")

	conversationID, err := strconv.Atoi(conversationIDStr)
	if err != nil {
		utils.WriteJSONError(w, http.StatusBadRequest, err)
		return
	}

	conversation, err := h.convService.GetConversation(r.Context(), int64(conversationID))
	if err != nil {
		utils.WriteJSONError(w, http.StatusBadRequest, err)
		return
	}

	utils.WriteJSON(w, http.StatusAccepted, conversation)
}

func (h *Handler) postConversationDM(w http.ResponseWriter, r *http.Request) {
	senderID := r.Context().Value(utils.CtxKeyUserID)

	var createConvReq pkg.CreateConversationDmRequest

	if err := json.NewDecoder(r.Body).Decode(&createConvReq); err != nil {
		utils.WriteJSONError(w, http.StatusBadRequest, err)
		return
	}

	conversation, err := h.convService.GetOrCreateDM(r.Context(), senderID.(int64), createConvReq.ParticipantID)
	if err != nil {
		utils.WriteJSONError(w, http.StatusBadRequest, err)
		return
	}

	utils.WriteJSON(w, http.StatusAccepted, conversation)
}

func (h *Handler) deleteConversation(w http.ResponseWriter, r *http.Request) {
}

func (h *Handler) postMessageInConversation(w http.ResponseWriter, r *http.Request) {
	senderID := r.Context().Value(utils.CtxKeyUserID)

	var msgReq pkg.SendMsgRequest

	if err := json.NewDecoder(r.Body).Decode(&msgReq); err != nil {
		utils.WriteJSONError(w, http.StatusBadRequest, err)
	}

	conversationIDStr := chi.URLParam(r, "conversation_id")

	conversationID, err := strconv.Atoi(conversationIDStr)
	if err != nil {
		utils.WriteJSONError(w, http.StatusBadRequest, err)
		return
	}

	if err := h.convService.PostConversationMsg(r.Context(), senderID.(int64), int64(conversationID), msgReq.Content); err != nil {
		utils.WriteJSONError(w, http.StatusBadRequest, err)
		return
	}

	utils.WriteJsonMsg(w, http.StatusAccepted, "Msg sent")
}
