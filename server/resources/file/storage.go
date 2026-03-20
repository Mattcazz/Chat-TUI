package file

import (
	"context"

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

func (s *FileStore) CreateFile(ctx context.Context, file *File) error {
	query := `INSERT INTO files (file_name, conversation_id, size, status, checksum, created_at) 
						VALUES ($1, $2, $3, $4, $5, $6) 
						RETURNING id`

	err := s.db.QueryRowContext(ctx, query, file.FileName, file.ConversationID, file.Size, file.Status, file.Checksum, file.CreatedAt).Scan(&file.ID)
	return err
}

func (s *FileStore) DeleteFile(ctx context.Context, fileID int64) error {
	query := `DELETE FROM files WHERE id = $1`
	_, err := s.db.ExecContext(ctx, query, fileID)
	return err
}

func (s *FileStore) InitUploadSession(ctx context.Context, uploadSession *UploadSession) error {
	query := `INSERT INTO upload_sessions (file_id, total_chunks, status, expired_at) 
						VALUES ($1, $2, $3, $4) 
						RETURNING id`

	err := s.db.QueryRowContext(ctx, query, uploadSession.FileID, uploadSession.TotalChunks, uploadSession.Status, uploadSession.ExpiresAt).Scan(&uploadSession.ID)
	return err
}

func (s *FileStore) DeleteUploadSession(ctx context.Context, sessionID int64) error {
	query := `DELETE FROM upload_sessions WHERE id = $1`
	_, err := s.db.ExecContext(ctx, query, sessionID)
	return err
}

func (s *FileStore) InsertFileChunk(ctx context.Context, fileChunk *FileChunk) error {
	query := `INSERT INTO file_chunks (index, session_id, created_at, checksum) 
						VALUES ($1, $2, $3, $4) 
						RETURNING id`

	err := s.db.QueryRowContext(ctx, query, fileChunk.Index, fileChunk.SessionID, fileChunk.CreatedAt, fileChunk.Checksum).Scan(&fileChunk.ID)
	return err
}

func (s *FileStore) DeleteFileChunksFromUploadSession(ctx context.Context, sessionID int64) error {
	query := `DELETE FROM file_chunks WHERE session_id = $1`
	_, err := s.db.ExecContext(ctx, query, sessionID)
	return err
}

func (s *FileStore) GetChunksCountForSession(ctx context.Context, sessionID int64) (int64, error) {
	query := `SELECT COUNT(*) FROM file_chunks WHERE session_id = $1`
	var count int64
	err := s.db.QueryRowContext(ctx, query, sessionID).Scan(&count)
	return count, err
}
