package chat

import (
	"context"
	"database/sql"

	"github.com/Mattcazz/Chat-TUI/server/db"
)

type ConversationStore struct {
	db db.DBTX
}

func NewConversationStore(db *sql.DB) *ConversationStore {
	return &ConversationStore{
		db: db,
	}
}

func (s *ConversationStore) WithTx(tx *sql.Tx) *ConversationStore {
	return &ConversationStore{
		db: tx,
	}
}

func (s *ConversationStore) CreateMessage(ctx context.Context, msg Message) error {
	return nil
}

func (s *ConversationStore) DeleteMessage(ctx context.Context, msgID int64) error {
	return nil
}

func (s *ConversationStore) GetMessage(context.Context, int64) *Message {
	return nil
}

func (s *ConversationStore) CreateConversation(context.Context, Conversation) error {
	return nil
}

func (s *ConversationStore) DeleteConversation(context.Context, int64) error {
	return nil
}

func (s *ConversationStore) EditConversation(context.Context, Conversation) error {
	return nil
}

func (s *ConversationStore) GetConversation(context.Context, int64) (*Conversation, error) {
	return nil, nil
}

func (s *ConversationStore) GetConversationHistory(ctx context.Context, converastionID, limit int64) (*[]Message, error) {
	return nil, nil
}
