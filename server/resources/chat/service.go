package chat

import "github.com/Mattcazz/Chat-TUI/server/db"

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
