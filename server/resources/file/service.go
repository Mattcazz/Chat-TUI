package file

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
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
		Checksum:       initFileReq.Checksum,
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

	dir := filepath.Join(string(TmpUploadsPath), fmt.Sprintf("session-%d", uploadSession.ID))
	err = os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return nil, err
	}

	if err := s.tx.Commit(tx); err != nil {
		os.RemoveAll(dir) // clean up the created directory if commit fails
		return nil, err
	}

	resp := &pkg.InitFileUploadResponse{
		SessionID: uploadSession.ID,
		FileID:    file.ID,
	}

	return resp, nil
}

func (s *Service) UploadFileChunk(ctx context.Context, uploadChunkReq *pkg.UploadFileChunkRequest) error {
	fileChunk := &FileChunk{
		Index:     uploadChunkReq.ChunkIndex,
		SessionID: uploadChunkReq.SessionID,
		CreatedAt: time.Now(),
		Checksum:  uploadChunkReq.Checksum,
	}

	path := filepath.Join(string(TmpUploadsPath), fmt.Sprintf("session-%d/chunk-%d.bin", uploadChunkReq.SessionID, uploadChunkReq.ChunkIndex))

	createdFile, err := os.Create(path)
	if err != nil {
		return err
	}

	defer createdFile.Close()

	n, err := createdFile.Write(uploadChunkReq.ChunkData)
	if err != nil {
		return err
	}

	if n != len(uploadChunkReq.ChunkData) {
		return fmt.Errorf("failed to write the entire chunk data to file, expected %d bytes but wrote %d bytes", len(uploadChunkReq.ChunkData), n)
	}
	return s.fileRepo.InsertFileChunk(ctx, fileChunk)
}

func (s *Service) FinalizeFileUpload(ctx context.Context, sessionID int64) error {
	session, err := s.fileRepo.GetUploadSession(ctx, sessionID)
	if err != nil {
		return err
	}
	if session == nil {
		return fmt.Errorf("upload session with id %d not found", sessionID)
	}

	if session.Status != FileSessionStatusUploading {
		return fmt.Errorf("cannot finalize upload session with status %s", session.Status)
	}

	chunksCount, err := s.fileRepo.GetChunksCountForSession(ctx, sessionID)
	if err != nil {
		return err
	}

	if chunksCount != session.TotalChunks {
		return fmt.Errorf("cannot finalize upload session because the number of uploaded chunks %d does not match the expected total chunks %d", chunksCount, session.TotalChunks)
	}

	file, err := s.fileRepo.GetFile(ctx, session.FileID)
	if err != nil {
		return err
	}

	filename := fmt.Sprintf("file-%d%s", session.FileID, file.Extension)
	finalPath := filepath.Join(string(FinalUploadsPath), filename)
	finalFile, err := os.Create(finalPath)
	if err != nil {
		return err
	}
	defer func() {
		finalFile.Close()
		if err != nil {
			os.Remove(finalPath) // clean up partial file on any error
		}
	}()

	for i := int64(0); i < session.TotalChunks; i++ {
		chunkPath := filepath.Join(string(TmpUploadsPath), fmt.Sprintf("session-%d/chunk-%d.bin", sessionID, i))

		chunk, err := os.Open(chunkPath)
		if err != nil {
			return err
		}

		if _, err := io.Copy(finalFile, chunk); err != nil {
			chunk.Close()
			return err
		}
		chunk.Close()
	}

	fileInfo, err := finalFile.Stat()
	if err != nil {
		return err
	}

	if fileInfo.Size() != file.Size {
		return fmt.Errorf("final assembled file size %d does not match expected file size %d", fileInfo.Size(), file.Size)
	}

	checkSum, err := pkg.CalculateFileChecksum(finalPath)
	if err != nil {
		return err
	}

	if checkSum != file.Checksum {
		return fmt.Errorf("final assembled file checksum does not match expected file checksum")
	}

	tx, err := s.tx.StartTx(ctx)
	if err != nil {
		return err
	}

	defer s.tx.RollBack(tx)
	if err := s.fileRepo.WithTx(tx).DeleteFileChunksFromUploadSession(ctx, sessionID); err != nil {
		return err
	}

	if err := s.fileRepo.WithTx(tx).UpdateFileStatus(ctx, session.FileID, FileStatusReady); err != nil {
		return err
	}

	if err := s.fileRepo.WithTx(tx).UpdateUploadSessionStatus(ctx, sessionID, FileSessionStatusCompleted); err != nil {
		return err
	}

	if err := s.tx.Commit(tx); err != nil {
		return err
	}

	sessionDir := filepath.Join(string(TmpUploadsPath), fmt.Sprintf("session-%d", sessionID))
	err = os.RemoveAll(sessionDir)
	if err != nil {
		return fmt.Errorf("failed to remove temporary upload session directory: %w", err)
	}

	return nil
}

func (s *Service) DeleteSessionChunks(ctx context.Context, sessionID int64) error {
	return nil
}
