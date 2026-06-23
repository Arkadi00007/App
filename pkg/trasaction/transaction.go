package transaction

import (
	"context"
	"database/sql"
)

type contextKey string

const txKey contextKey = "tx"

type Manager struct {
	db *sql.DB
}

func NewManager(db *sql.DB) *Manager {
	return &Manager{db: db}
}

// выполняет функцию внутри транзакции
// если fn вернёт ошибку — транзакция откатывается
// если fn вернёт nil — транзакция коммитится
func (m *Manager) WithTx(ctx context.Context, fn func(ctx context.Context) error) error {
	tx, err := m.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	// кладём tx в context чтобы repository мог его достать
	txCtx := context.WithValue(ctx, txKey, tx)

	if err := fn(txCtx); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

// repository вызывает эту функцию чтобы получить либо tx, либо обычный db
func getExecutor(ctx context.Context, db *sql.DB) interface {
	ExecContext(context.Context, string, ...any) (sql.Result, error)
	QueryRowContext(context.Context, string, ...any) *sql.Row
} {
	if tx, ok := ctx.Value(txKey).(*sql.Tx); ok {
		return tx
	}
	return db
}
