package tx

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"

	"github.com/jmoiron/sqlx"
)

type SqlxTransaction struct {
	db *sqlx.DB
}

func NewSqlxTransaction(
	db *sqlx.DB,
) *SqlxTransaction {
	return &SqlxTransaction{
		db: db,
	}
}

func (t *SqlxTransaction) Do(
	ctx context.Context,
	fn func(context.Context) error,
) error {
	tx, err := t.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("tx: failed to begin transaction: %w", err)
	}
	defer func() {
		if err := tx.Rollback(); err != nil && err != sql.ErrTxDone {
			slog.Error("tx: failed to rollback transaction", "error", err)
		}
	}()

	ctx = context.WithValue(ctx, ContextKey, tx)
	if err := fn(ctx); err != nil {
		return fmt.Errorf("tx: failed to execute function: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("tx: failed to commit transaction: %w", err)
	}

	return nil
}
