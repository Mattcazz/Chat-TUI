package file

import (
	"context"
	"fmt"
	"io"
	"log"
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
	log.Printf("Service.InitFileUpload: Initiating file upload - filename: %s, size: %d bytes, chunks: %d, conversation ID: %d",
		initFileReq.FileName, initFileReq.TotalSize, initFileReq.TotalChunks, initFileReq.ConversationID)

	fileName := SanitizeFileName(initFileReq.FileName)
	fileExtenstion := GetFileExtension(fileName)

	log.Printf("Service.InitFileUpload: Sanitized filename: %s, extension: %s", fileName, fileExtenstion)

	file := &File{
		FileName:       fileName,
		Extension:      fileExtenstion,
		ConversationID: initFileReq.ConversationID,
		UploaderID:     initFileReq.UploaderID,
		Size:           initFileReq.TotalSize,
		Checksum:       initFileReq.Checksum,
		Status:         FileStatusUploading,
		StoragePath:    "", // will be set upon finalization
		CreatedAt:      time.Now(),
	}

	log.Printf("Service.InitFileUpload: Starting transaction for file creation")
	tx, err := s.tx.StartTx(ctx)
	if err != nil {
		log.Printf("Service.InitFileUpload: Failed to start transaction: %v", err)
		return nil, err
	}

	defer s.tx.RollBack(tx)

	if err := s.fileRepo.WithTx(tx).CreateFile(ctx, file); err != nil {
		log.Printf("Service.InitFileUpload: Failed to create file record: %v", err)
		return nil, err
	}

	log.Printf("Service.InitFileUpload: File record created with ID %d", file.ID)

	uploadSession := &UploadSession{
		FileID:      file.ID,
		TotalChunks: initFileReq.TotalChunks,
		Status:      FileSessionStatusUploading,
		ExpiresAt:   time.Now().Add(time.Duration(TimeToExpireUploadSession) * time.Second),
	}

	if err := s.fileRepo.WithTx(tx).InitUploadSession(ctx, uploadSession); err != nil {
		log.Printf("Service.InitFileUpload: Failed to create upload session: %v", err)
		return nil, err
	}

	log.Printf("Service.InitFileUpload: Upload session created with ID %d", uploadSession.ID)

	dir := filepath.Join(string(TmpUploadsPath), fmt.Sprintf("session-%d", uploadSession.ID))
	log.Printf("Service.InitFileUpload: Creating temporary directory: %s", dir)
	err = os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		log.Printf("Service.InitFileUpload: Failed to create temporary directory: %v", err)
		return nil, err
	}

	log.Printf("Service.InitFileUpload: Committing transaction")
	if err := s.tx.Commit(tx); err != nil {
		log.Printf("Service.InitFileUpload: Failed to commit transaction, cleaning up directory: %v", err)
		os.RemoveAll(dir) // clean up the created directory if commit fails
		return nil, err
	}

	resp := &pkg.InitFileUploadResponse{
		SessionID: uploadSession.ID,
		FileID:    file.ID,
	}

	log.Printf("Service.InitFileUpload: Successfully initialized file upload - session ID: %d, file ID: %d", uploadSession.ID, file.ID)
	return resp, nil
}

func (s *Service) UploadFileChunk(ctx context.Context, sessionID int64, uploadChunkReq *pkg.UploadFileChunkRequest) error {
	log.Printf("Service.UploadFileChunk: Uploading chunk %d for session ID %d, data size: %d bytes",
		uploadChunkReq.ChunkIndex, sessionID, len(uploadChunkReq.ChunkData))

	fileChunk := &FileChunk{
		Index:     uploadChunkReq.ChunkIndex,
		SessionID: sessionID,
		CreatedAt: time.Now(),
		Checksum:  uploadChunkReq.Checksum,
	}

	path := filepath.Join(string(TmpUploadsPath), fmt.Sprintf("session-%d/chunk-%d.bin", sessionID, uploadChunkReq.ChunkIndex))
	log.Printf("Service.UploadFileChunk: Writing chunk to file: %s", path)

	createdFile, err := os.Create(path)
	if err != nil {
		log.Printf("Service.UploadFileChunk: Failed to create chunk file: %v", err)
		return err
	}

	defer createdFile.Close()

	n, err := createdFile.Write(uploadChunkReq.ChunkData)
	if err != nil {
		log.Printf("Service.UploadFileChunk: Failed to write chunk data: %v", err)
		return err
	}

	if n != len(uploadChunkReq.ChunkData) {
		os.Remove(path)
		return fmt.Errorf("failed to write the entire chunk data to file, expected %d bytes but wrote %d bytes", len(uploadChunkReq.ChunkData), n)
	}

	log.Printf("Service.UploadFileChunk: Inserting chunk record into database")
	if err := s.fileRepo.InsertFileChunk(ctx, fileChunk); err != nil {
		os.Remove(path)
		log.Printf("Service.UploadFileChunk: Failed to insert chunk record: %v", err)
		return err
	}

	log.Printf("Service.UploadFileChunk: Successfully uploaded chunk %d for session ID %d", uploadChunkReq.ChunkIndex, sessionID)
	return nil
}

func (s *Service) FinalizeFileUpload(ctx context.Context, sessionID int64) error {
	log.Printf("Service.FinalizeFileUpload: Starting finalization for session ID %d", sessionID)

	session, err := s.fileRepo.GetUploadSession(ctx, sessionID)
	if err != nil {
		log.Printf("Service.FinalizeFileUpload: Failed to retrieve session: %v", err)
		return err
	}
	if session == nil {
		log.Printf("Service.FinalizeFileUpload: Session ID %d not found", sessionID)
		return fmt.Errorf("upload session with id %d not found", sessionID)
	}

	log.Printf("Service.FinalizeFileUpload: Session retrieved with status: %s", session.Status)

	if session.Status != FileSessionStatusUploading {
		return fmt.Errorf("cannot finalize upload session with status %s", session.Status)
	}

	chunksCount, err := s.fileRepo.GetChunksCountForSession(ctx, sessionID)
	if err != nil {
		log.Printf("Service.FinalizeFileUpload: Failed to get chunks count: %v", err)
		return err
	}

	log.Printf("Service.FinalizeFileUpload: Validating chunks - received: %d, expected: %d", chunksCount, session.TotalChunks)

	if chunksCount != session.TotalChunks {
		log.Printf("Service.FinalizeFileUpload: Chunk count mismatch for session %d", sessionID)
		return fmt.Errorf("cannot finalize upload session because the number of uploaded chunks %d does not match the expected total chunks %d", chunksCount, session.TotalChunks)
	}

	file, err := s.fileRepo.GetFile(ctx, session.FileID)
	if err != nil {
		log.Printf("Service.FinalizeFileUpload: Failed to retrieve file metadata: %v", err)
		return err
	}

	log.Printf("Service.FinalizeFileUpload: Retrieved file metadata - filename: %s, extension: %s, size: %d bytes",
		file.FileName, file.Extension, file.Size)

	filename := fmt.Sprintf("file-%d%s", session.FileID, file.Extension)
	finalPath := filepath.Join(string(FinalUploadsPath), filename)

	log.Printf("Service.FinalizeFileUpload: Creating final file at: %s", finalPath)
	finalFile, err := os.Create(finalPath)
	if err != nil {
		log.Printf("Service.FinalizeFileUpload: Failed to create final file: %v", err)
		return err
	}
	defer func() {
		finalFile.Close()
		if err != nil {
			os.Remove(finalPath) // clean up partial file on any error
		}
	}()

	log.Printf("Service.FinalizeFileUpload: Starting file assembly from %d chunks", session.TotalChunks)
	if err := assembleFile(finalFile, finalPath, file, session); err != nil {
		return err
	}

	log.Printf("Service.FinalizeFileUpload: File assembled successfully, starting database finalization")
	if err := s.finalizeFileUploadOnDB(ctx, session, finalPath); err != nil {
		return err
	}

	sessionDir := filepath.Join(string(TmpUploadsPath), fmt.Sprintf("session-%d", sessionID))
	log.Printf("Service.FinalizeFileUpload: Cleaning up temporary directory: %s", sessionDir)
	err = os.RemoveAll(sessionDir)
	if err != nil {
		return fmt.Errorf("failed to remove temporary upload session directory: %w", err)
	}

	log.Printf("Service.FinalizeFileUpload: Successfully finalized file upload for session ID %d", sessionID)
	return nil
}

func assembleFile(finalFile *os.File, finalPath string, file *File, session *UploadSession) error {
	log.Printf("assembleFile: Assembling file from %d chunks for session ID %d, expected size: %d bytes",
		session.TotalChunks, session.ID, file.Size)

	for i := int64(0); i < session.TotalChunks; i++ {
		chunkPath := filepath.Join(string(TmpUploadsPath), fmt.Sprintf("session-%d/chunk-%d.bin", session.ID, i))

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

	log.Printf("assembleFile: Validating file size - actual: %d bytes, expected: %d bytes", fileInfo.Size(), file.Size)

	if fileInfo.Size() != file.Size {
		return fmt.Errorf("final assembled file size %d does not match expected file size %d", fileInfo.Size(), file.Size)
	}

	log.Printf("assembleFile: Calculating and validating file checksum")
	checkSum, err := pkg.CalculateFileChecksum(finalPath)
	if err != nil {
		return err
	}

	if checkSum != file.Checksum {
		return fmt.Errorf("final assembled file checksum does not match expected file checksum")
	}

	log.Printf("assembleFile: File assembly completed successfully - size validated, checksum validated")
	return nil
}

func (s *Service) finalizeFileUploadOnDB(ctx context.Context, session *UploadSession, finalPath string) error {
	log.Printf("Service.finalizeFileUploadOnDB: Starting database finalization for session ID %d, file ID %d",
		session.ID, session.FileID)

	tx, err := s.tx.StartTx(ctx)
	if err != nil {
		return err
	}

	defer s.tx.RollBack(tx)

	log.Printf("Service.finalizeFileUploadOnDB: Deleting chunk records for session %d", session.ID)
	if err := s.fileRepo.WithTx(tx).DeleteFileChunksFromUploadSession(ctx, session.ID); err != nil {
		return err
	}

	log.Printf("Service.finalizeFileUploadOnDB: Updating file ID %d status to ready", session.FileID)
	if err := s.fileRepo.WithTx(tx).UpdateFileStatusAndPath(ctx, session.FileID, FileStatusReady, finalPath); err != nil {
		return err
	}

	log.Printf("Service.finalizeFileUploadOnDB: Updating session ID %d status to completed", session.ID)
	if err := s.fileRepo.WithTx(tx).UpdateUploadSessionStatus(ctx, session.ID, FileSessionStatusCompleted); err != nil {
		return err
	}

	log.Printf("Service.finalizeFileUploadOnDB: Committing transaction")
	if err := s.tx.Commit(tx); err != nil {
		return err
	}

	log.Printf("Service.finalizeFileUploadOnDB: Database finalization completed successfully")
	return nil
}

func (s *Service) DeleteSessionChunks(ctx context.Context, sessionID int64) error {
	log.Printf("Service.DeleteSessionChunks: Function not yet implemented for session ID %d", sessionID)
	return nil
}
