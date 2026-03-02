package chat

import (
	"context"
	"database/sql"

	"github.com/Mattcazz/Chat-TUI/pkg"
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

func (s *ConversationStore) AddParticipantToConversation(ctx context.Context, conversationID, participantID int64, role ParticipantRole) error {
	query := `INSERT INTO conversation_participants (convconversation_id, user_id, role) VALUES (1$, 2$, 3$)`

	_, err := s.db.ExecContext(ctx, query, conversationID, participantID, role)
	return err
}

func (s *ConversationStore) CreateMessage(ctx context.Context, msg *Message) error {
	query := `INSERT INTO messages (conversation_id, sender_id, content, created_at) VALUES (1$, 2$, 3$, 4$)`

	_, err := s.db.ExecContext(ctx, query, msg.ConversationID, msg.SenderID, msg.Content, msg.CreatedAt)
	if err != nil {
		return err
	}

	return nil
}

func (s *ConversationStore) DeleteMessage(ctx context.Context, msgID int64) error {
	query := `DELETE FROM messages WHERE id = 1$`

	_, err := s.db.ExecContext(ctx, query, msgID)
	return err
}

func (s *ConversationStore) GetMessage(context.Context, int64) *Message {
	return nil
}

func (s *ConversationStore) CreateConversation(ctx context.Context, conversation *Conversation) error {
	query := `INSERT INTO conversations (created_at) VALUES (1$) RETURNING id`

	row, err := s.db.QueryContext(ctx, query, conversation.CreatedAt)
	if err != nil {
		return err
	}

	return row.Scan(conversation.ID)
}

func (s *ConversationStore) DeleteConversation(context.Context, int64) error {
	return nil
}

func (s *ConversationStore) EditConversation(context.Context, *Conversation) error {
	return nil
}

func (s *ConversationStore) GetConversation(ctx context.Context, id, limit int64) (*pkg.ConversationResponse, error) {
	query := `SELECT m.content, u.username, m.created_at
						FROM messages m 
						LEFT JOIN user u
						ON m.user_id = u.id
						WHERE m.conversation_id = $1
						ORDER BY m.created_at DESC
						LIMIT $2`

	rows, err := s.db.QueryContext(ctx, query, id, limit)
	if err != nil {
		return nil, err
	}

	conversation := &pkg.ConversationResponse{
		ID:       id,
		Messages: []pkg.MsgResponse{},
	}

	for rows.Next() {
		var msg pkg.MsgResponse
		if err := rows.Scan(&msg.Content, &msg.UserName, &msg.CreatedAt); err != nil {
			return nil, err
		}
	}

	return conversation, nil
}

func (s *ConversationStore) GetConversationHistory(ctx context.Context, converastionID, limit int64) (*pkg.ConversationResponse, error) {
	return nil, nil
}

// TODO: nickname and see if i need to return the username of the other person.
func (s *ConversationStore) GetConversationDM(ctx context.Context, firstID, secondID, limit int64) (*pkg.ConversationResponse, error) {
	query := `WITH conversation AS (
						SELECT conversation_id 
						FROM conversation_participants
						WHERE user_id = $1
						INTERSECT
						SELECT conversation_id
						FROM conversation_participants
						WHERE user_id = $2)
						SELECT m.conversation_id, m.content, u.username, m.created_at
						FROM messages m
						LEFT JOIN users u ON m.user_id = u.id
						WHERE m.conversation_id IN (SELECT conversation_id FROM conversation)
						ORDER BY m.created_at DESC
						LIMIT $3`

	conversation := &pkg.ConversationResponse{
		Messages: []pkg.MsgResponse{},
	}
	rows, err := s.db.QueryContext(ctx, query, firstID, secondID, limit)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var msg pkg.MsgResponse
		if err := rows.Scan(&conversation.ID, &msg.Content, &msg.UserName, &msg.CreatedAt); err != nil {
			return nil, err
		}

		conversation.Messages = append(conversation.Messages, msg)
	}

	return conversation, nil
}
