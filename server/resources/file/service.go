package file

import (
	"context"
	"time"

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

func (s *Service) InitFileUpload(ctx context.Context, initFileReq *pkg.InitFileUploadRequest) (*pkg.InitFileUploadResponse, error) {
	file := &File{
		FileName:       initFileReq.FileName,
		ConversationID: initFileReq.ConversationID,
		Size:           initFileReq.TotalSize,
		Checksum:       "", // TODO: what to do here? client should send checksum of the file or we calculate it on the server after receiving all chunks?
		Status:         FileStatusUploading,
		CreatedAt:      time.Now(),
	}

	tx, err := s.tx.StartTx(ctx)
	if err != nil {
		return nil, err
	}

	defer s.tx.RollBack(tx)

	if err := s.fileRepo.WithTx(tx).CreateFile(ctx, file); err != nil {
		return nil, err
	}

	uploadSession := &UploadSession{
		FileID:      file.ID,
		TotalChunks: initFileReq.TotalChunks,
		Status:      FileSessionStatusUploading,
		ExpiresAt:   time.Now().Add(time.Duration(TimeToExpireUploadSession) * time.Second),
	}

	if err := s.fileRepo.WithTx(tx).InitUploadSession(ctx, uploadSession); err != nil {
		return nil, err
	}

	resp := &pkg.InitFileUploadResponse{
		SessionID: uploadSession.ID,
		FileID:    file.ID,
	}

	if err := s.tx.Commit(tx); err != nil {
		return nil, err
	}

	return resp, nil
}

func (s *Service) UploadFileChunk(ctx context.Context, uploadChunkReq *pkg.UploadFileChunkRequest) error {
	fileChunk := &FileChunk{
		Index:     uploadChunkReq.ChunkIndex,
		SessionID: uploadChunkReq.SessionID,
		CreatedAt: time.Now(),
		Checksum:  "", // TODO: same as above, client should send checksum of the chunk or we calculate it on the server after receiving the chunk?
	}

	return s.fileRepo.InsertFileChunk(ctx, fileChunk)
}

func (s *Service) FinalizeFileUpload(ctx context.Context, sessionID int64) error {
	return nil
}

func (s *Service) DeleteSessionChunks(ctx context.Context, sessionID int64) error {
	return nil
}
