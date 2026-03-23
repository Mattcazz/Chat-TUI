package file

import (
	"context"
	"database/sql"
	"time"
)

type FileRepository interface {
	WithTx(tx *sql.Tx) FileRepository
	CreateFile(ctx context.Context, file *File) error
	DeleteFile(ctx context.Context, fileID int64) error
	InitUploadSession(ctx context.Context, uploadSession *UploadSession) error
	DeleteUploadSession(ctx context.Context, sessionID int64) error
	InsertFileChunk(ctx context.Context, fileChunk *FileChunk) error
	DeleteFileChunksFromUploadSession(ctx context.Context, sessionID int64) error
	GetChunksCountForSession(ctx context.Context, sessionID int64) (int64, error)
	UpdateFileStatus(ctx context.Context, fileID int64, status FileStatus) error
	UpdateUploadSessionStatus(ctx context.Context, sessionID int64, status UploadSessionStatus) error
	GetUploadSession(ctx context.Context, sessionID int64) (*UploadSession, error)
}

type File struct {
	ID             int64      `json:"id"`
	FileName       string     `json:"file_name"`
	ConversationID int64      `json:"conversation_id"`
	Size           int64      `json:"size"`
	Status         FileStatus `json:"status"`
	Checksum       string     `json:"checksum"`
	CreatedAt      time.Time  `json:"created_at"`
}

type UploadSession struct {
	ID          int64               `json:"id"`
	FileID      int64               `json:"file_id"`
	TotalChunks int64               `json:"total_chunks"`
	Status      UploadSessionStatus `json:"status"`
	ExpiresAt   time.Time           `json:"expires_at"`
}

type FileChunk struct {
	ID        int64     `json:"id"`
	Index     int64     `json:"index"`
	SessionID int64     `json:"session_id"`
	CreatedAt time.Time `json:"created_at"`
	Checksum  string    `json:"checksum"`
}
