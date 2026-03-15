package chat

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Mattcazz/Chat-TUI/pkg"
	"github.com/Mattcazz/Chat-TUI/server/utils"
	"github.com/go-chi/chi/v5"
)

// ---- Helpers ----

// injectChiParam sets a chi URL param on the request context — needed because
// chi.URLParam reads from its own context key, not the URL string.
func injectChiParam(r *http.Request, key, value string) *http.Request {
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add(key, value)
	return r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))
}

// withUserID injects a fake authenticated user ID the same way the JWT middleware does.
func withUserID(r *http.Request, id int64) *http.Request {
	return r.WithContext(context.WithValue(r.Context(), utils.CtxKeyUserID, id))
}

func newTestHandler(repo ConversationRepository) *Handler {
	broker := NewBroker()
	svc := &Service{conversationRepo: repo, broker: broker}
	return NewHandler(svc, broker)
}

// ---- getConversation ----

func TestHandler_GetConversation_OK(t *testing.T) {
	repo := &mockConversationRepo{
		getConversationFn: func(_ context.Context, id, _ int64) (*pkg.ConversationResponse, error) {
			return &pkg.ConversationResponse{ID: id, Messages: []pkg.MsgResponse{}}, nil
		},
	}
	h := newTestHandler(repo)

	req := httptest.NewRequest(http.MethodGet, "/conversation/7", nil)
	req = injectChiParam(req, "conversation_id", "7")
	w := httptest.NewRecorder()

	h.getConversation(w, req)

	if w.Code != http.StatusAccepted {
		t.Errorf("status: got %d, want %d", w.Code, http.StatusAccepted)
	}

	var resp pkg.ConversationResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("could not decode response: %v", err)
	}
	if resp.ID != 7 {
		t.Errorf("ID: got %d, want 7", resp.ID)
	}
}

func TestHandler_GetConversation_InvalidID(t *testing.T) {
	h := newTestHandler(&mockConversationRepo{})

	req := httptest.NewRequest(http.MethodGet, "/conversation/abc", nil)
	req = injectChiParam(req, "conversation_id", "abc")
	w := httptest.NewRecorder()

	h.getConversation(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("status: got %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func TestHandler_GetConversation_ServiceError(t *testing.T) {
	repo := &mockConversationRepo{
		getConversationFn: func(_ context.Context, _ int64, _ int64) (*pkg.ConversationResponse, error) {
			return nil, errors.New("not found")
		},
	}
	h := newTestHandler(repo)

	req := httptest.NewRequest(http.MethodGet, "/conversation/1", nil)
	req = injectChiParam(req, "conversation_id", "1")
	w := httptest.NewRecorder()

	h.getConversation(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("status: got %d, want %d", w.Code, http.StatusBadRequest)
	}
}

// ---- postConversationDM ----

func TestHandler_PostConversationDM_OK(t *testing.T) {
	repo := &mockConversationRepo{
		getConversationDMFn: func(_ context.Context, _, _, _ int64) (*pkg.ConversationResponse, error) {
			return &pkg.ConversationResponse{ID: 3, Messages: []pkg.MsgResponse{}}, nil
		},
	}
	h := newTestHandler(repo)

	body, _ := json.Marshal(pkg.CreateConversationDmRequest{ParticipantID: 2})
	req := httptest.NewRequest(http.MethodPost, "/conversation", bytes.NewReader(body))
	req = withUserID(req, 1)
	w := httptest.NewRecorder()

	h.postConversationDM(w, req)

	if w.Code != http.StatusAccepted {
		t.Errorf("status: got %d, want %d", w.Code, http.StatusAccepted)
	}
}

func TestHandler_PostConversationDM_InvalidBody(t *testing.T) {
	h := newTestHandler(&mockConversationRepo{})

	req := httptest.NewRequest(http.MethodPost, "/conversation", bytes.NewReader([]byte("not json")))
	req = withUserID(req, 1)
	w := httptest.NewRecorder()

	h.postConversationDM(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("status: got %d, want %d", w.Code, http.StatusBadRequest)
	}
}

// ---- postMessageInConversation ----

func TestHandler_PostMessageInConversation_OK(t *testing.T) {
	h := newTestHandler(&mockConversationRepo{})

	body, _ := json.Marshal(pkg.SendMsgRequest{Content: "hello"})
	req := httptest.NewRequest(http.MethodPost, "/conversation/1/message", bytes.NewReader(body))
	req = withUserID(req, 1)
	req = injectChiParam(req, "conversation_id", "1")
	w := httptest.NewRecorder()

	h.postMessageInConversation(w, req)

	if w.Code != http.StatusAccepted {
		t.Errorf("status: got %d, want %d", w.Code, http.StatusAccepted)
	}
}

func TestHandler_PostMessageInConversation_InvalidConvID(t *testing.T) {
	h := newTestHandler(&mockConversationRepo{})

	body, _ := json.Marshal(pkg.SendMsgRequest{Content: "hello"})
	req := httptest.NewRequest(http.MethodPost, "/conversation/abc/message", bytes.NewReader(body))
	req = withUserID(req, 1)
	req = injectChiParam(req, "conversation_id", "abc")
	w := httptest.NewRecorder()

	h.postMessageInConversation(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("status: got %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func TestHandler_PostMessageInConversation_InvalidBody(t *testing.T) {
	h := newTestHandler(&mockConversationRepo{})

	req := httptest.NewRequest(http.MethodPost, "/conversation/1/message", bytes.NewReader([]byte("bad")))
	req = withUserID(req, 1)
	req = injectChiParam(req, "conversation_id", "1")
	w := httptest.NewRecorder()

	h.postMessageInConversation(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("status: got %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func TestHandler_PostMessageInConversation_ServiceError(t *testing.T) {
	repo := &mockConversationRepo{
		createMessageFn: func(_ context.Context, _ *Message) (*pkg.MsgResponse, error) {
			return nil, errors.New("failed")
		},
	}
	h := newTestHandler(repo)

	body, _ := json.Marshal(pkg.SendMsgRequest{Content: "hello"})
	req := httptest.NewRequest(http.MethodPost, "/conversation/1/message", bytes.NewReader(body))
	req = withUserID(req, 1)
	req = injectChiParam(req, "conversation_id", "1")
	w := httptest.NewRecorder()

	h.postMessageInConversation(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("status: got %d, want %d", w.Code, http.StatusBadRequest)
	}
}
