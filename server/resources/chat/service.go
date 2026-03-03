package chat

import (
	"context"
	"log"
	"time"

	"github.com/Mattcazz/Chat-TUI/pkg"
	"github.com/Mattcazz/Chat-TUI/server/db"
)

type Service struct {
	conversationRepo ConversationRepository
	tx               *db.TxManager
}

func NewService(conversationRepo ConversationRepository, tx *db.TxManager) *Service {
	return &Service{
		conversationRepo: conversationRepo,
		tx:               tx,
	}
}

func (s *Service) PostConversationMsg(ctx context.Context, senderID, convID int64, content string) error {
	msg := &Message{
		SenderID:       senderID,
		ConversationID: convID,
		Content:        content,
		CreatedAt:      time.Now(),
	}
	return s.conversationRepo.CreateMessage(ctx, msg)
}

// user wants to start conversation with a contact
func (s *Service) GetOrCreateDM(ctx context.Context, creatorID, participantID int64) (*pkg.ConversationResponse, error) {
	conversation, err := s.conversationRepo.GetConversationDM(ctx, creatorID, participantID, 20)
	if err != nil {
		return nil, err
	}

	if conversation == nil || conversation.ID == 0 {
		var participants []int64
		participants = append(participants, participantID)
		return s.createChat(ctx, creatorID, participants)
	}

	return conversation, nil
}

func (s *Service) GetConversation(ctx context.Context, conversationID int64) (*pkg.ConversationResponse, error) {
	return s.conversationRepo.GetConversation(ctx, conversationID, 20) // TODO: How many messages?
}

func (s *Service) createChat(ctx context.Context, creatorID int64, otherParticipants []int64) (*pkg.ConversationResponse, error) {
	conversation := &Conversation{
		CreatedAt: time.Now(),
	}

	tx, err := s.tx.StartTx(ctx)
	if err != nil {
		return nil, err
	}

	defer tx.Rollback()

	if err := s.conversationRepo.WithTx(tx).CreateConversation(ctx, conversation); err != nil {
		return nil, err
	} // we will get the id of the conversation after this point.

	log.Println("Created conversation ", conversation.ID)
	// then we add the creator as admin. TODO: decide if this logic makes sense (admin role for creator)
	if err := s.conversationRepo.WithTx(tx).AddParticipantToConversation(ctx, conversation.ID, creatorID, RoleAdmin); err != nil {
		return nil, err
	}

	var partID int64
	for i := range len(otherParticipants) { // TODO: Check logic of WithTx to see if the performance is super bad in a loop
		partID = otherParticipants[i]
		if err := s.conversationRepo.WithTx(tx).AddParticipantToConversation(ctx, conversation.ID, partID, RoleAdmin); err != nil {
			return nil, err
		}
	}
	tx.Commit()

	return s.conversationRepo.GetConversation(ctx, conversation.ID, 20) // TODO: decide how many messages
}
