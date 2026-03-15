package chat

import (
	"context"
	"database/sql"
	"time"

	"github.com/Mattcazz/Chat-TUI/pkg"
)

type ConversationRepository interface {
	WithTx(*sql.Tx) *ConversationStore
	AddParticipantToConversation(ctx context.Context, conversationID, participantID int64, role ParticipantRole) error
	CreateConversation(context.Context, *Conversation) error
	DeleteConversation(context.Context, int64) error
	EditConversation(context.Context, *Conversation) error
	GetConversation(context.Context, int64, int64) (*pkg.ConversationResponse, error)
	DeleteMessage(context.Context, int64) error
	GetMessage(context.Context, int64) (*Message, error)
	CreateMessage(context.Context, *Message) (*pkg.MsgResponse, error)
	GetConversationDM(ctx context.Context, firstID, secondID, limit int64) (*pkg.ConversationResponse, error)
}

type Message struct {
	SenderID       int64     `json:"sender_id"`
	Content        string    `json:"content"`
	ConversationID int64     `json:"conversation_id"`
	CreatedAt      time.Time `json:"created_at"`
}

type Conversation struct {
	ID        int64     `json:"id"`
	LastMsg   string    `json:"last_message"`
	LastMsgAt time.Time `json:"last_message_at"`
	CreatedAt time.Time `json:"created_at"`
}

type ParticipantRole string

const (
	RoleAdmin  ParticipantRole = "admin"
	RoleMember ParticipantRole = "member"
)
