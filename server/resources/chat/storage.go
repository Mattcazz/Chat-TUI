package chat

import (
	"context"
	"database/sql"
	"fmt"
	"log"

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
	log.Printf("ConversationStore.AddParticipantToConversation: Adding participant ID %d to conversation ID %d with role %s", participantID, conversationID, role)
	query := `INSERT INTO conversation_participants (conversation_id, user_id, role) VALUES ($1, $2, $3)`

	_, err := s.db.ExecContext(ctx, query, conversationID, participantID, role)
	if err != nil {
		log.Printf("ConversationStore.AddParticipantToConversation: Failed to add participant ID %d to conversation ID %d: %v", participantID, conversationID, err)
		return err
	}

	log.Printf("ConversationStore.AddParticipantToConversation: Successfully added participant ID %d to conversation ID %d", participantID, conversationID)
	return nil
}

func (s *ConversationStore) CreateMessage(ctx context.Context, msg *Message) (*pkg.MsgResponse, error) {
	log.Printf("ConversationStore.CreateMessage: Creating message in conversation ID %d from sender ID %d", msg.ConversationID, msg.SenderID)
	query := `WITH inserted as (
						INSERT INTO messages (conversation_id, sender_id, content, created_at, message_type)
						VALUES ($1, $2, $3, $4, $5)
						RETURNING *),
						updated_conversation as (
						UPDATE conversations SET last_message_at = $4, last_message_preview = $3
						WHERE id = $1
						RETURNING id)

						SELECT u.username, inserted.content, inserted.created_at, inserted.message_type
						FROM inserted JOIN users u ON inserted.sender_id = u.id`

	rows, err := s.db.QueryContext(ctx, query, msg.ConversationID, msg.SenderID, msg.Content, msg.CreatedAt, msg.Type)
	if err != nil {
		log.Printf("ConversationStore.CreateMessage: Failed to create message in conversation ID %d: %v", msg.ConversationID, err)
		return nil, err
	}
	defer rows.Close()

	if rows.Next() {
		response := &pkg.MsgResponse{}
		if err := rows.Scan(&response.UserName, &response.Content, &response.CreatedAt, &response.Type); err != nil {
			log.Printf("ConversationStore.CreateMessage: Failed to scan message response in conversation ID %d: %v", msg.ConversationID, err)
			return nil, err
		}
		log.Printf("ConversationStore.CreateMessage: Successfully created message in conversation ID %d", msg.ConversationID)
		return response, nil
	}

	log.Printf("ConversationStore.CreateMessage: Failed to insert message in conversation ID %d", msg.ConversationID)
	return nil, fmt.Errorf("failed to insert message")
}

func (s *ConversationStore) DeleteMessage(ctx context.Context, msgID int64) error {
	log.Printf("ConversationStore.DeleteMessage: Deleting message with ID %d", msgID)
	query := `DELETE FROM messages WHERE id = $1`

	_, err := s.db.ExecContext(ctx, query, msgID)
	if err != nil {
		log.Printf("ConversationStore.DeleteMessage: Failed to delete message with ID %d: %v", msgID, err)
		return err
	}

	log.Printf("ConversationStore.DeleteMessage: Successfully deleted message with ID %d", msgID)
	return nil
}

func (s *ConversationStore) GetMessage(ctx context.Context, id int64) (*Message, error) {
	log.Printf("ConversationStore.GetMessage: Retrieving message with ID %d", id)
	query := `SELECT sender_id, content, conversation_id, created_at, message_type FROM messages WHERE id = $1`

	var msg Message
	err := s.db.QueryRowContext(ctx, query, id).Scan(&msg.SenderID, &msg.Content, &msg.ConversationID, &msg.CreatedAt, &msg.Type)
	if err != nil {
		log.Printf("ConversationStore.GetMessage: Failed to retrieve message with ID %d: %v", id, err)
		return nil, err
	}

	log.Printf("ConversationStore.GetMessage: Successfully retrieved message with ID %d", id)
	return &msg, nil
}

func (s *ConversationStore) CreateConversation(ctx context.Context, conversation *Conversation) error {
	log.Printf("ConversationStore.CreateConversation: Creating new conversation")
	query := `INSERT INTO conversations (created_at) VALUES ($1) RETURNING id`

	err := s.db.QueryRowContext(ctx, query, conversation.CreatedAt).Scan(&conversation.ID)
	if err != nil {
		log.Printf("ConversationStore.CreateConversation: Failed to create conversation: %v", err)
		return fmt.Errorf("failed to insert conversation: %w", err)
	}

	log.Printf("ConversationStore.CreateConversation: Successfully created conversation with ID %d", conversation.ID)
	return nil
}

func (s *ConversationStore) DeleteConversation(ctx context.Context, id int64) error {
	log.Printf("ConversationStore.DeleteConversation: Deleting conversation with ID %d", id)
	query := `DELETE FROM conversations WHERE id = $1`

	_, err := s.db.ExecContext(ctx, query, id)
	if err != nil {
		log.Printf("ConversationStore.DeleteConversation: Failed to delete conversation with ID %d: %v", id, err)
		return err
	}

	log.Printf("ConversationStore.DeleteConversation: Successfully deleted conversation with ID %d", id)
	return nil
}

func (s *ConversationStore) EditConversation(ctx context.Context, conversation *Conversation) error {
	log.Printf("ConversationStore.EditConversation: Updating conversation with ID %d", conversation.ID)
	query := `UPDATE conversations SET last_message_at = $1, last_message_preview = $2
						WHERE id = $3`

	err := s.db.QueryRowContext(ctx, query, conversation.LastMsgAt, conversation.LastMsg, conversation.ID).Err()
	if err != nil {
		log.Printf("ConversationStore.EditConversation: Failed to update conversation with ID %d: %v", conversation.ID, err)
		return err
	}

	log.Printf("ConversationStore.EditConversation: Successfully updated conversation with ID %d", conversation.ID)
	return nil
}

func (s *ConversationStore) GetConversation(ctx context.Context, id, limit int64) (*pkg.ConversationResponse, error) {
	log.Printf("ConversationStore.GetConversation: Retrieving conversation ID %d with limit %d", id, limit)
	query := `SELECT m.content, u.username, m.created_at, m.message_type	
						FROM messages m
						LEFT JOIN users u
						ON m.sender_id = u.id
						WHERE m.conversation_id = $1
						ORDER BY m.created_at DESC
						LIMIT $2`

	rows, err := s.db.QueryContext(ctx, query, id, limit)
	if err != nil {
		log.Printf("ConversationStore.GetConversation: Failed to query messages for conversation ID %d: %v", id, err)
		return nil, err
	}

	conversation := &pkg.ConversationResponse{
		ID:       id,
		Messages: []pkg.MsgResponse{},
	}

	for rows.Next() {
		var msg pkg.MsgResponse
		if err := rows.Scan(&msg.Content, &msg.UserName, &msg.CreatedAt, &msg.Type); err != nil {
			log.Printf("ConversationStore.GetConversation: Failed to scan message row for conversation ID %d: %v", id, err)
			return nil, err
		}
		conversation.Messages = append(conversation.Messages, msg)
	}

	if err := rows.Err(); err != nil {
		log.Printf("ConversationStore.GetConversation: Error iterating message rows for conversation ID %d: %v", id, err)
		return nil, err
	}

	log.Printf("ConversationStore.GetConversation: Successfully retrieved %d messages for conversation ID %d", len(conversation.Messages), id)
	return conversation, nil
}

// TODO: nickname and see if i need to return the username of the other person.
func (s *ConversationStore) GetConversationDM(ctx context.Context, firstID, secondID, limit int64) (*pkg.ConversationResponse, error) {
	log.Printf("ConversationStore.GetConversationDM: Retrieving DM conversation between user ID %d and user ID %d with limit %d", firstID, secondID, limit)
	query := `WITH conversation AS (
						SELECT conversation_id
						FROM conversation_participants
						WHERE user_id = $1
						INTERSECT
						SELECT conversation_id
						FROM conversation_participants
						WHERE user_id = $2)
						SELECT m.conversation_id, m.content, u.username, m.created_at, m.message_type
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
		log.Printf("ConversationStore.GetConversationDM: Failed to query DM conversation between user ID %d and user ID %d: %v", firstID, secondID, err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var msg pkg.MsgResponse
		if err := rows.Scan(&conversation.ID, &msg.Content, &msg.UserName, &msg.CreatedAt, &msg.Type); err != nil {
			log.Printf("ConversationStore.GetConversationDM: Failed to scan message row for DM between user ID %d and user ID %d: %v", firstID, secondID, err)
			return nil, err
		}

		conversation.Messages = append(conversation.Messages, msg)
	}

	if err := rows.Err(); err != nil {
		log.Printf("ConversationStore.GetConversationDM: Error iterating message rows for DM between user ID %d and user ID %d: %v", firstID, secondID, err)
		return nil, err
	}

	log.Printf("ConversationStore.GetConversationDM: Successfully retrieved %d messages for DM conversation between user ID %d and user ID %d", len(conversation.Messages), firstID, secondID)
	return conversation, nil
}
