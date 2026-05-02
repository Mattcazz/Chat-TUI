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
	broker           *Broker
}

func NewService(conversationRepo ConversationRepository, tx *db.TxManager, broker *Broker) *Service {
	return &Service{
		conversationRepo: conversationRepo,
		tx:               tx,
		broker:           broker,
	}
}

func (s *Service) PostConversationMsg(ctx context.Context, senderID, convID int64, content string, msgType pkg.MsgType) error {
	log.Printf("Service.PostConversationMsg: Posting message to conversation ID %d from user ID %d", convID, senderID)
	msg := &Message{
		SenderID:       senderID,
		ConversationID: convID,
		Content:        content,
		Type:           msgType,
		CreatedAt:      time.Now(),
	}

	msgResp, err := s.conversationRepo.CreateMessage(ctx, msg)
	if err != nil {
		log.Printf("Service.PostConversationMsg: Failed to create message in conversation ID %d: %v", convID, err)
		return err
	}

	s.broker.Publish(convID, *msgResp)
	log.Printf("Service.PostConversationMsg: Successfully published message to conversation ID %d", convID)

	return nil
}

// user wants to start conversation with a contact
func (s *Service) GetOrCreateDM(ctx context.Context, creatorID, participantID int64) (*pkg.ConversationResponse, error) {
	log.Printf("Service.GetOrCreateDM: Getting or creating DM between user ID %d and user ID %d", creatorID, participantID)
	conversation, err := s.conversationRepo.GetConversationDM(ctx, creatorID, participantID, 20)
	if err != nil {
		log.Printf("Service.GetOrCreateDM: Failed to get DM conversation between user ID %d and user ID %d: %v", creatorID, participantID, err)
		return nil, err
	}

	if conversation == nil || conversation.ID == 0 {
		log.Printf("Service.GetOrCreateDM: DM conversation not found, creating new conversation between user ID %d and user ID %d", creatorID, participantID)
		var participants []int64
		participants = append(participants, participantID)
		return s.createChat(ctx, creatorID, participants)
	}

	log.Printf("Service.GetOrCreateDM: Successfully retrieved existing DM conversation with ID %d", conversation.ID)
	return conversation, nil
}

func (s *Service) GetConversation(ctx context.Context, conversationID int64) (*pkg.ConversationResponse, error) {
	log.Printf("Service.GetConversation: Retrieving conversation ID %d", conversationID)
	conversation, err := s.conversationRepo.GetConversation(ctx, conversationID, 20) // TODO: How many messages?
	if err != nil {
		log.Printf("Service.GetConversation: Failed to retrieve conversation ID %d: %v", conversationID, err)
		return nil, err
	}
	log.Printf("Service.GetConversation: Successfully retrieved conversation ID %d with %d messages", conversationID, len(conversation.Messages))
	return conversation, nil
}

func (s *Service) createChat(ctx context.Context, creatorID int64, otherParticipants []int64) (*pkg.ConversationResponse, error) {
	log.Printf("Service.createChat: Creating new conversation with creator ID %d and %d other participants", creatorID, len(otherParticipants))
	conversation := &Conversation{
		CreatedAt: time.Now(),
	}

	tx, err := s.tx.StartTx(ctx)
	if err != nil {
		log.Printf("Service.createChat: Failed to start transaction: %v", err)
		return nil, err
	}

	defer tx.Rollback()

	if err := s.conversationRepo.WithTx(tx).CreateConversation(ctx, conversation); err != nil {
		log.Printf("Service.createChat: Failed to create conversation: %v", err)
		return nil, err
	} // we will get the id of the conversation after this point.

	log.Printf("Service.createChat: Created conversation with ID %d", conversation.ID)
	// then we add the creator as admin. TODO: decide if this logic makes sense (admin role for creator)
	if err := s.conversationRepo.WithTx(tx).AddParticipantToConversation(ctx, conversation.ID, creatorID, RoleAdmin); err != nil {
		log.Printf("Service.createChat: Failed to add creator to conversation ID %d: %v", conversation.ID, err)
		return nil, err
	}

	var partID int64
	for i := range len(otherParticipants) { // TODO: Check logic of WithTx to see if the performance is super bad in a loop
		partID = otherParticipants[i]
		if err := s.conversationRepo.WithTx(tx).AddParticipantToConversation(ctx, conversation.ID, partID, RoleAdmin); err != nil {
			log.Printf("Service.createChat: Failed to add participant to conversation ID %d: %v", conversation.ID, err)
			return nil, err
		}
	}
	if err := tx.Commit(); err != nil {
		log.Printf("Service.createChat: Failed to commit transaction for conversation ID %d: %v", conversation.ID, err)
		return nil, err
	}

	log.Printf("Service.createChat: Successfully created conversation ID %d with all participants", conversation.ID)
	return s.conversationRepo.GetConversation(ctx, conversation.ID, 20) // TODO: decide how many messages
}
