package file

import (
	"context"
	"database/sql"
	"log"

	"github.com/Mattcazz/Chat-TUI/server/db"
)

type FileStore struct {
	db db.DBTX
}

func NewFileStore(db db.DBTX) *FileStore {
	return &FileStore{
		db: db,
	}
}

func (s *FileStore) WithTx(tx *sql.Tx) *FileStore {
	return &FileStore{
		db: tx,
	}
}

func (s *FileStore) GetFile(ctx context.Context, fileID int64) (*File, error) {
	log.Printf("FileStore.GetFile: Retrieving file with ID %d", fileID)
	query := `SELECT id, file_name, extension, conversation_id, uploader_id, size, status, checksum, created_at FROM files WHERE id = $1`

	var file File
	err := s.db.QueryRowContext(ctx, query, fileID).Scan(
		&file.ID,
		&file.FileName,
		&file.Extension,
		&file.ConversationID,
		&file.UploaderID,
		&file.Size,
		&file.Status,
		&file.Checksum,
		&file.CreatedAt,
	)
	if err != nil {
		log.Printf("FileStore.GetFile: Failed to retrieve file with ID %d: %v", fileID, err)
		return nil, err
	}

	log.Printf("FileStore.GetFile: Successfully retrieved file with ID %d, filename %s", file.ID, file.FileName)
	return &file, nil
}

func (s *FileStore) CreateFile(ctx context.Context, file *File) error {
	log.Printf("FileStore.CreateFile: Creating file with name %s, size %d bytes, conversation ID %d", file.FileName, file.Size, file.ConversationID)
	query := `INSERT INTO files (file_name, extension, conversation_id, uploader_id, size, status, checksum, created_at) 
						VALUES ($1, $2, $3, $4, $5, $6, $7, $8) 
						RETURNING id`

	err := s.db.QueryRowContext(ctx, query, file.FileName, file.Extension, file.ConversationID, file.UploaderID, file.Size, file.Status, file.Checksum, file.CreatedAt).Scan(&file.ID)
	if err != nil {
		log.Printf("FileStore.CreateFile: Failed to create file %s: %v", file.FileName, err)
		return err
	}

	log.Printf("FileStore.CreateFile: Successfully created file with ID %d", file.ID)
	return nil
}

func (s *FileStore) DeleteFile(ctx context.Context, fileID int64) error {
	log.Printf("FileStore.DeleteFile: Deleting file with ID %d", fileID)
	query := `DELETE FROM files WHERE id = $1`
	_, err := s.db.ExecContext(ctx, query, fileID)
	if err != nil {
		log.Printf("FileStore.DeleteFile: Failed to delete file with ID %d: %v", fileID, err)
		return err
	}

	log.Printf("FileStore.DeleteFile: Successfully deleted file with ID %d", fileID)
	return nil
}

func (s *FileStore) GetUploadSession(ctx context.Context, sessionID int64) (*UploadSession, error) {
	log.Printf("FileStore.GetUploadSession: Retrieving upload session with ID %d", sessionID)
	query := `SELECT id, file_id, total_chunks, status, expired_at FROM upload_sessions WHERE id = $1`
	row := s.db.QueryRowContext(ctx, query, sessionID)

	var session UploadSession
	err := row.Scan(&session.ID, &session.FileID, &session.TotalChunks, &session.Status, &session.ExpiresAt)
	if err != nil {
		log.Printf("FileStore.GetUploadSession: Failed to retrieve session with ID %d: %v", sessionID, err)
		return nil, err
	}

	log.Printf("FileStore.GetUploadSession: Successfully retrieved session with ID %d, status %s, total chunks %d", session.ID, session.Status, session.TotalChunks)
	return &session, nil
}

func (s *FileStore) InitUploadSession(ctx context.Context, uploadSession *UploadSession) error {
	log.Printf("FileStore.InitUploadSession: Creating upload session for file ID %d with %d total chunks", uploadSession.FileID, uploadSession.TotalChunks)
	query := `INSERT INTO upload_sessions (file_id, total_chunks, status, expired_at) 
						VALUES ($1, $2, $3, $4) 
						RETURNING id`

	err := s.db.QueryRowContext(ctx, query, uploadSession.FileID, uploadSession.TotalChunks, uploadSession.Status, uploadSession.ExpiresAt).Scan(&uploadSession.ID)
	if err != nil {
		log.Printf("FileStore.InitUploadSession: Failed to create session for file ID %d: %v", uploadSession.FileID, err)
		return err
	}

	log.Printf("FileStore.InitUploadSession: Successfully created session with ID %d", uploadSession.ID)
	return nil
}

func (s *FileStore) DeleteUploadSession(ctx context.Context, sessionID int64) error {
	log.Printf("FileStore.DeleteUploadSession: Deleting upload session with ID %d", sessionID)
	query := `DELETE FROM upload_sessions WHERE id = $1`
	_, err := s.db.ExecContext(ctx, query, sessionID)
	if err != nil {
		log.Printf("FileStore.DeleteUploadSession: Failed to delete session with ID %d: %v", sessionID, err)
		return err
	}

	log.Printf("FileStore.DeleteUploadSession: Successfully deleted session with ID %d", sessionID)
	return nil
}

func (s *FileStore) InsertFileChunk(ctx context.Context, fileChunk *FileChunk) error {
	log.Printf("FileStore.InsertFileChunk: Inserting chunk %d for session ID %d", fileChunk.Index, fileChunk.SessionID)
	query := `INSERT INTO file_chunks (index, session_id, created_at, checksum) 
						VALUES ($1, $2, $3, $4) 
						RETURNING id`

	err := s.db.QueryRowContext(ctx, query, fileChunk.Index, fileChunk.SessionID, fileChunk.CreatedAt, fileChunk.Checksum).Scan(&fileChunk.ID)
	if err != nil {
		log.Printf("FileStore.InsertFileChunk: Failed to insert chunk %d for session %d: %v", fileChunk.Index, fileChunk.SessionID, err)
		return err
	}

	log.Printf("FileStore.InsertFileChunk: Successfully inserted chunk %d with ID %d for session %d", fileChunk.Index, fileChunk.ID, fileChunk.SessionID)
	return nil
}

func (s *FileStore) DeleteFileChunksFromUploadSession(ctx context.Context, sessionID int64) error {
	log.Printf("FileStore.DeleteFileChunksFromUploadSession: Deleting chunks for session ID %d", sessionID)
	query := `DELETE FROM file_chunks WHERE session_id = $1`
	_, err := s.db.ExecContext(ctx, query, sessionID)
	if err != nil {
		log.Printf("FileStore.DeleteFileChunksFromUploadSession: Failed to delete chunks for session %d: %v", sessionID, err)
		return err
	}

	log.Printf("FileStore.DeleteFileChunksFromUploadSession: Successfully deleted chunks for session ID %d", sessionID)
	return nil
}

func (s *FileStore) GetChunksCountForSession(ctx context.Context, sessionID int64) (int64, error) {
	log.Printf("FileStore.GetChunksCountForSession: Counting chunks for session ID %d", sessionID)
	query := `SELECT COUNT(*) FROM file_chunks WHERE session_id = $1`
	var count int64
	err := s.db.QueryRowContext(ctx, query, sessionID).Scan(&count)
	if err != nil {
		log.Printf("FileStore.GetChunksCountForSession: Failed to count chunks for session %d: %v", sessionID, err)
		return 0, err
	}

	log.Printf("FileStore.GetChunksCountForSession: Found %d chunks for session ID %d", count, sessionID)
	return count, nil
}

func (s *FileStore) UpdateFileStatus(ctx context.Context, fileID int64, status FileStatus) error {
	log.Printf("FileStore.UpdateFileStatus: Updating file ID %d status to %s", fileID, status)
	query := `UPDATE files SET status = $1 WHERE id = $2`

	_, err := s.db.ExecContext(ctx, query, status, fileID)
	if err != nil {
		log.Printf("FileStore.UpdateFileStatus: Failed to update file ID %d status: %v", fileID, err)
		return err
	}

	log.Printf("FileStore.UpdateFileStatus: Successfully updated file ID %d to status %s", fileID, status)
	return nil
}

func (s *FileStore) UpdateUploadSessionStatus(ctx context.Context, sessionID int64, status UploadSessionStatus) error {
	log.Printf("FileStore.UpdateUploadSessionStatus: Updating session ID %d status to %s", sessionID, status)
	query := `UPDATE upload_sessions SET status = $1 WHERE id = $2`

	_, err := s.db.ExecContext(ctx, query, status, sessionID)
	if err != nil {
		log.Printf("FileStore.UpdateUploadSessionStatus: Failed to update session ID %d status: %v", sessionID, err)
		return err
	}

	log.Printf("FileStore.UpdateUploadSessionStatus: Successfully updated session ID %d to status %s", sessionID, status)
	return nil
}
