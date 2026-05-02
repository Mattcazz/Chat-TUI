package chat

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/Mattcazz/Chat-TUI/pkg"
)

// ---- Mock repository ----

type mockConversationRepo struct {
	// Control return values per method
	getConversationFn    func(ctx context.Context, id, limit int64) (*pkg.ConversationResponse, error)
	getConversationDMFn  func(ctx context.Context, firstID, secondID, limit int64) (*pkg.ConversationResponse, error)
	createConversationFn func(ctx context.Context, c *Conversation) error
	addParticipantFn     func(ctx context.Context, convID, participantID int64, role ParticipantRole) error
	createMessageFn      func(ctx context.Context, msg *Message) (*pkg.MsgResponse, error)
}

func (m *mockConversationRepo) WithTx(_ *sql.Tx) *ConversationStore { return nil }

func (m *mockConversationRepo) GetConversation(ctx context.Context, id, limit int64) (*pkg.ConversationResponse, error) {
	if m.getConversationFn != nil {
		return m.getConversationFn(ctx, id, limit)
	}
	return &pkg.ConversationResponse{ID: id}, nil
}

func (m *mockConversationRepo) GetConversationDM(ctx context.Context, firstID, secondID, limit int64) (*pkg.ConversationResponse, error) {
	if m.getConversationDMFn != nil {
		return m.getConversationDMFn(ctx, firstID, secondID, limit)
	}
	return &pkg.ConversationResponse{}, nil
}

func (m *mockConversationRepo) CreateConversation(ctx context.Context, c *Conversation) error {
	if m.createConversationFn != nil {
		return m.createConversationFn(ctx, c)
	}
	c.ID = 99 // simulate DB returning an ID
	return nil
}

func (m *mockConversationRepo) AddParticipantToConversation(ctx context.Context, convID, participantID int64, role ParticipantRole) error {
	if m.addParticipantFn != nil {
		return m.addParticipantFn(ctx, convID, participantID, role)
	}
	return nil
}

func (m *mockConversationRepo) CreateMessage(ctx context.Context, msg *Message) (*pkg.MsgResponse, error) {
	if m.createMessageFn != nil {
		return m.createMessageFn(ctx, msg)
	}
	return &pkg.MsgResponse{UserName: "alice", Content: msg.Content}, nil
}

func (m *mockConversationRepo) DeleteConversation(ctx context.Context, id int64) error { return nil }
func (m *mockConversationRepo) EditConversation(ctx context.Context, c *Conversation) error {
	return nil
}
func (m *mockConversationRepo) DeleteMessage(ctx context.Context, id int64) error { return nil }
func (m *mockConversationRepo) GetMessage(ctx context.Context, id int64) (*Message, error) {
	return nil, nil
}

// ---- Helpers ----

func newTestService(repo ConversationRepository) *Service {
	broker := NewBroker()
	// TxManager is nil — tests that hit createChat will need a real one or we skip those paths
	return &Service{
		conversationRepo: repo,
		broker:           broker,
	}
}

// ---- Service tests ----

func TestService_GetConversation(t *testing.T) {
	repo := &mockConversationRepo{
		getConversationFn: func(_ context.Context, id, _ int64) (*pkg.ConversationResponse, error) {
			return &pkg.ConversationResponse{ID: id, Messages: []pkg.MsgResponse{}}, nil
		},
	}
	svc := newTestService(repo)

	conv, err := svc.GetConversation(context.Background(), 42)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if conv.ID != 42 {
		t.Errorf("got ID %d, want 42", conv.ID)
	}
}

func TestService_GetConversation_Error(t *testing.T) {
	repo := &mockConversationRepo{
		getConversationFn: func(_ context.Context, _ int64, _ int64) (*pkg.ConversationResponse, error) {
			return nil, errors.New("db error")
		},
	}
	svc := newTestService(repo)

	_, err := svc.GetConversation(context.Background(), 1)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestService_PostConversationMsg_PublishesToBroker(t *testing.T) {
	repo := &mockConversationRepo{}
	broker := NewBroker()
	svc := &Service{conversationRepo: repo, broker: broker}

	ch := broker.Subscribe(10)

	if err := svc.PostConversationMsg(context.Background(), 1, 10, "hi there", pkg.MsgTypeText); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	select {
	case msg := <-ch:
		if msg.Content != "hi there" {
			t.Errorf("got content %q, want %q", msg.Content, "hi there")
		}
	case <-time.After(time.Second):
		t.Error("timed out waiting for message from broker")
	}
}

func TestService_PostConversationMsg_RepoError(t *testing.T) {
	repo := &mockConversationRepo{
		createMessageFn: func(_ context.Context, _ *Message) (*pkg.MsgResponse, error) {
			return nil, errors.New("insert failed")
		},
	}
	svc := newTestService(repo)

	err := svc.PostConversationMsg(context.Background(), 1, 1, "hello", pkg.MsgTypeText)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestService_GetOrCreateDM_ExistingConversation(t *testing.T) {
	existing := &pkg.ConversationResponse{
		ID:       5,
		Messages: []pkg.MsgResponse{{UserName: "bob", Content: "hey", Type: pkg.MsgTypeText}},
	}
	repo := &mockConversationRepo{
		getConversationDMFn: func(_ context.Context, _, _, _ int64) (*pkg.ConversationResponse, error) {
			return existing, nil
		},
	}
	svc := newTestService(repo)

	conv, err := svc.GetOrCreateDM(context.Background(), 1, 2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if conv.ID != 5 {
		t.Errorf("got ID %d, want 5", conv.ID)
	}
}

func TestService_GetOrCreateDM_DMRepoError(t *testing.T) {
	repo := &mockConversationRepo{
		getConversationDMFn: func(_ context.Context, _, _, _ int64) (*pkg.ConversationResponse, error) {
			return nil, errors.New("db down")
		},
	}
	svc := newTestService(repo)

	_, err := svc.GetOrCreateDM(context.Background(), 1, 2)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
