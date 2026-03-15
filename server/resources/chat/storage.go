package chat

import (
	"context"
	"database/sql"
	"fmt"

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
	query := `INSERT INTO conversation_participants (conversation_id, user_id, role) VALUES ($1, $2, $3)`

	_, err := s.db.ExecContext(ctx, query, conversationID, participantID, role)
	return err
}

func (s *ConversationStore) CreateMessage(ctx context.Context, msg *Message) (*pkg.MsgResponse, error) {
	query := `WITH inserted as (
						INSERT INTO messages (conversation_id, sender_id, content, created_at) 
						VALUES ($1, $2, $3, $4)
						RETURNING *)

						SELECT u.username, inserted.content, inserted.created_at
						FROM inserted JOIN users u ON inserted.sender_id = u.id`

	rows, err := s.db.QueryContext(ctx, query, msg.ConversationID, msg.SenderID, msg.Content, msg.CreatedAt)
	if err != nil {
		return nil, err
	}

	if rows.Next() {
		response := &pkg.MsgResponse{}
		rows.Scan(&response.UserName, &response.Content, &response.CreatedAt)
		return response, nil
	} else {
		return nil, fmt.Errorf("failed to insert message")
	}
}

func (s *ConversationStore) DeleteMessage(ctx context.Context, msgID int64) error {
	query := `DELETE FROM messages WHERE id = $1`

	_, err := s.db.ExecContext(ctx, query, msgID)
	return err
}

func (s *ConversationStore) GetMessage(ctx context.Context, id int64) (*Message, error) {
	query := `SELECT sender_id, content, conversation_id, created_at FROM messages WHERE id = $1`

	var msg Message
	err := s.db.QueryRowContext(ctx, query, id).Scan(&msg.SenderID, &msg.Content, &msg.ConversationID, &msg.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &msg, nil
}

func (s *ConversationStore) CreateConversation(ctx context.Context, conversation *Conversation) error {
	query := `INSERT INTO conversations (created_at) VALUES ($1) RETURNING id`

	err := s.db.QueryRowContext(ctx, query, conversation.CreatedAt).Scan(&conversation.ID)
	if err != nil {
		return fmt.Errorf("failed to insert conversation: %w", err)
	}

	return nil
}

func (s *ConversationStore) DeleteConversation(ctx context.Context, id int64) error {
	query := `DELETE FROM conversations WHERE id = $1`

	_, err := s.db.ExecContext(ctx, query, id)

	return err
}

func (s *ConversationStore) EditConversation(ctx context.Context, conversation *Conversation) error {
	query := `UPDATE conversations SET last_message_at = $1, last_message_preview = $2 
						WHERE id = $3`

	return s.db.QueryRowContext(ctx, query, conversation.LastMsgAt, conversation.LastMsg, conversation.ID).Err()
}

func (s *ConversationStore) GetConversation(ctx context.Context, id, limit int64) (*pkg.ConversationResponse, error) {
	query := `SELECT m.content, u.username, m.created_at
						FROM messages m 
						LEFT JOIN users u
						ON m.sender_id = u.id
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
		conversation.Messages = append(conversation.Messages, msg)
	}

	return conversation, nil
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
						LEFT JOIN users u ON m.sender_id = u.id
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
