package file

import "github.com/Mattcazz/Chat-TUI/server/db"

type Service struct {
	fileRepo FileRepository
	tx       *db.TxManager
}

func NewService(fr FileRepository, tx *db.TxManager) *Service {
	return &Service{
		fileRepo: fr,
		tx:       tx,
	}
}
