package file

type FileRepository interface {
	CreateFile(file *File) error
	DeleteFile(fileID int64) error
	InitUploadSession(uploadSession *UploadSession) error
	DeleteUploadSession(sessionID int64) error
	InsertFileChunk(fileChunk *FileChunk) error
	DeleteFileChunk(chunkID int64) error
}

type File struct {
	ID            int64  `json:"id"`
	FileName      string `json:"file_name"`
	ConvesationID int64  `json:"conversation_id"`
	Size          int64  `json:"size"`
	Status        string `json:"status"`
	Checksum      string `json:"checksum"`
	CreatedAt     int64  `json:"created_at"`
}

type UploadSession struct {
	ID          int64  `json:"id"`
	FileID      int64  `json:"file_id"`
	TotalChunks int64  `json:"total_chunks"`
	Status      string `json:"status"`
	ExpiredAt   int64  `json:"expired_at"`
}

type FileChunk struct {
	ID        int64  `json:"id"`
	Index     int64  `json:"index"`
	SessionID int64  `json:"session_id"`
	CretaedAt int64  `json:"created_at"`
	Checksum  string `json:"checksum"`
}
