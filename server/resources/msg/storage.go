package msg

import "database/sql"

type MsgStore struct {
	db *sql.DB
}

func NewMsgStore(db *sql.DB) *MsgStore {
	return &MsgStore{
		db: db,
	}
}
