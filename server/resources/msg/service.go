package msg

import "github.com/Mattcazz/Chat-TUI/server/db"

type Service struct {
	msgRepo MsgRepository
	tx      *db.TxManager
}

func NewService(msgRepo MsgRepository, tx *db.TxManager) *Service {
	return &Service{
		msgRepo: msgRepo,
		tx:      tx,
	}
}
