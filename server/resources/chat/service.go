package chat

import (
	"context"

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

func (s *Service) postConversationMsg(ctx context.Context, sender_id, conv_id int64, content string) error {
	return nil
}
