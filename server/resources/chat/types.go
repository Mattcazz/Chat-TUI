package chat

import (
	"context"
	"database/sql"
	"time"
)

type ConversationRepository interface {
	WithTx(*sql.Tx) *ConversationStore
	CreateConversation(context.Context, Conversation) error
	DeleteConversation(context.Context, int64) error
	EditConversation(context.Context, Conversation) error
	GetConversation(context.Context, int64) (*Conversation, error)
	GetConversationHistory(ctx context.Context, converastionID, limit int64) (*[]Message, error)
	DeleteMessage(context.Context, int64) error
	GetMessage(context.Context, int64) *Message
	CreateMessage(context.Context, Message) error
}

type Message struct {
	SenderID       int64
	Content        string
	ConversationID int64
	CreatedAt      time.Time
}

type Conversation struct {
	ID        int64
	LastMsg   string
	LastMsgAt time.Time
	CreatedAt time.Time
}
