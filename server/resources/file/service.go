package file

import (
	"context"

	"github.com/Mattcazz/Chat-TUI/pkg"
	"github.com/Mattcazz/Chat-TUI/server/db"
)

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

func (s *Service) InitFileUpload(ctx context.Context, initFileReq *pkg.InitFileUploadRequest) error {
	return nil
}

func (s *Service) UploadFileChunk(ctx context.Context, uploadChunkReq *pkg.UploadFileChunkRequest) error {
	return nil
}

func (s *Service) FinalizeFileUpload(ctx context.Context, sessionID int64) error {
	return nil
}

func (s *Service) DeleteSessionChunks(ctx context.Context, sessionID int64) error {
	return nil
}
