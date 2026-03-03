package db

import (
	"context"
	"database/sql"
)

type DBTX interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
}

type TxManager struct {
	db *sql.DB
}

func NewTxManager(db *sql.DB) *TxManager {
	return &TxManager{
		db: db,
	}
}

func (t *TxManager) StartTx(ctx context.Context) (*sql.Tx, error) {
	return t.db.BeginTx(ctx, nil)
}

func (t *TxManager) Commit(tx *sql.Tx) error {
	return tx.Commit()
}

func (t *TxManager) RollBack(tx *sql.Tx) error {
	return tx.Rollback()
}
