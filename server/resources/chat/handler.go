package chat

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/Mattcazz/Chat-TUI/pkg"
	"github.com/Mattcazz/Chat-TUI/server/resources/middleware"
	"github.com/Mattcazz/Chat-TUI/server/utils"
	"github.com/go-chi/chi/v5"
)

type Handler struct {
	convService *Service
	broker      *Broker
}

func NewHandler(s *Service, broker *Broker) *Handler {
	return &Handler{
		convService: s,
		broker:      broker,
	}
}

func (h *Handler) RegisterRoutes(r *chi.Mux) {
	r.Route("/conversation", func(r chi.Router) {
		r.Post("/", middleware.JWTAuth(h.postConversationDM))
		r.Get("/{conversation_id}", middleware.JWTAuth(h.getConversation))
		r.Delete("/{conversation_id}", middleware.JWTAuth(h.deleteConversation))
		r.Post("/{conversation_id}/message", middleware.JWTAuth(h.postMessageInConversation))
		r.Get("/{conversation_id}/stream", middleware.JWTAuth(h.streamConversationMessages))
	})
}

func (h *Handler) getConversation(w http.ResponseWriter, r *http.Request) {
	conversationIDStr := chi.URLParam(r, "conversation_id")
	log.Printf("Handler.getConversation: Retrieving conversation with ID %s", conversationIDStr)

	conversationID, err := strconv.Atoi(conversationIDStr)
	if err != nil {
		log.Printf("Handler.getConversation: Invalid conversation ID format: %v", err)
		utils.WriteJSONError(w, http.StatusBadRequest, err)
		return
	}

	conversation, err := h.convService.GetConversation(r.Context(), int64(conversationID))
	if err != nil {
		log.Printf("Handler.getConversation: Failed to get conversation ID %d: %v", conversationID, err)
		utils.WriteJSONError(w, http.StatusBadRequest, err)
		return
	}

	log.Printf("Handler.getConversation: Successfully retrieved conversation ID %d", conversationID)
	utils.WriteJSON(w, http.StatusAccepted, conversation)
}

func (h *Handler) postConversationDM(w http.ResponseWriter, r *http.Request) {
	senderID := r.Context().Value(utils.CtxKeyUserID)
	log.Printf("Handler.postConversationDM: User ID %d requesting to create or get DM conversation", senderID)

	var createConvReq pkg.CreateConversationDmRequest

	if err := json.NewDecoder(r.Body).Decode(&createConvReq); err != nil {
		log.Printf("Handler.postConversationDM: Failed to decode request body: %v", err)
		utils.WriteJSONError(w, http.StatusBadRequest, err)
		return
	}

	log.Printf("Handler.postConversationDM: Getting or creating DM with participant ID %d", createConvReq.ParticipantID)
	conversation, err := h.convService.GetOrCreateDM(r.Context(), senderID.(int64), createConvReq.ParticipantID)
	if err != nil {
		log.Printf("Handler.postConversationDM: Failed to get or create DM conversation: %v", err)
		utils.WriteJSONError(w, http.StatusBadRequest, err)
		return
	}

	log.Printf("Handler.postConversationDM: Successfully retrieved or created conversation ID %d", conversation.ID)
	utils.WriteJSON(w, http.StatusAccepted, conversation)
}

func (h *Handler) deleteConversation(w http.ResponseWriter, r *http.Request) {
}

func (h *Handler) postMessageInConversation(w http.ResponseWriter, r *http.Request) {
	senderID := r.Context().Value(utils.CtxKeyUserID)
	log.Printf("Handler.postMessageInConversation: User ID %d posting message", senderID)

	var msgReq pkg.SendMsgRequest

	if err := json.NewDecoder(r.Body).Decode(&msgReq); err != nil {
		log.Printf("Handler.postMessageInConversation: Failed to decode request body: %v", err)
		utils.WriteJSONError(w, http.StatusBadRequest, err)
		return
	}

	conversationIDStr := chi.URLParam(r, "conversation_id")
	log.Printf("Handler.postMessageInConversation: Posting to conversation ID %s", conversationIDStr)

	conversationID, err := strconv.Atoi(conversationIDStr)
	if err != nil {
		log.Printf("Handler.postMessageInConversation: Invalid conversation ID format: %v", err)
		utils.WriteJSONError(w, http.StatusBadRequest, err)
		return
	}

	if err := h.convService.PostConversationMsg(r.Context(), senderID.(int64), int64(conversationID), msgReq.Content); err != nil {
		log.Printf("Handler.postMessageInConversation: Failed to post message to conversation ID %d: %v", conversationID, err)
		utils.WriteJSONError(w, http.StatusBadRequest, err)
		return
	}

	log.Printf("Handler.postMessageInConversation: Successfully posted message to conversation ID %d", conversationID)
	utils.WriteJsonMsg(w, http.StatusAccepted, "Msg sent")
}

func (h *Handler) streamConversationMessages(w http.ResponseWriter, r *http.Request) {
	conversationIDStr := chi.URLParam(r, "conversation_id")
	log.Printf("Handler.streamConversationMessages: Starting message stream for conversation ID %s", conversationIDStr)

	conversationID, err := strconv.Atoi(conversationIDStr)
	if err != nil {
		log.Printf("Handler.streamConversationMessages: Invalid conversation ID format: %v", err)
		utils.WriteJSONError(w, http.StatusBadRequest, err)
		return
	}

	flusher, ok := w.(http.Flusher) // if the ResponseWriter supports flushing (which is required for streaming)
	if !ok {
		log.Printf("Handler.streamConversationMessages: Streaming not supported by response writer for conversation ID %d", conversationID)
		http.Error(w, "streaming not supported", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	ch := h.broker.Subscribe(int64(conversationID))
	log.Printf("Handler.streamConversationMessages: Client subscribed to conversation ID %d", conversationID)
	defer h.broker.Unsubscribe(int64(conversationID), ch)

	for {
		select {
		case msg := <-ch:
			utils.WriteJSON(w, http.StatusOK, msg)
			flusher.Flush() // send whatever is in the buffer to the client immediately
		case <-r.Context().Done():
			log.Printf("Handler.streamConversationMessages: Client disconnected from conversation ID %d", conversationID)
			return // client disconnected
		}
	}
}
