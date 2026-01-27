package file

import "database/sql"

type FileStore struct {
	db *sql.DB
}

func NewFileStore(db *sql.DB) *FileStore {
	return &FileStore{
		db: db,
	}
}
