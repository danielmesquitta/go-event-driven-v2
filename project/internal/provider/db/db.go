package db

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"tickets/internal/pkg/tx"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var timeout = 10 * time.Second

type DB struct {
	*sqlx.DB
}

func NewSQLXDB() *sqlx.DB {
	db, err := sqlx.Open("postgres", os.Getenv("POSTGRES_URL"))
	if err != nil {
		panic(err)
	}
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		panic(err)
	}
	return db
}

func NewDB(db *sqlx.DB) *DB {
	return &DB{
		DB: db,
	}
}

func (d *DB) InitializeSchema(ctx context.Context) error {
	_, err := d.DB.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS tickets (
			id UUID PRIMARY KEY,
			price_amount   NUMERIC(10, 2) NOT NULL,
			price_currency CHAR(3) NOT NULL,
			customer_email VARCHAR(255) NOT NULL
		);

		CREATE TABLE IF NOT EXISTS shows (
			id UUID PRIMARY KEY,
			dead_nation_id UUID NOT NULL,
			number_of_tickets INT NOT NULL,
			start_time TIMESTAMP NOT NULL,
			title VARCHAR(255) NOT NULL,
			venue VARCHAR(255) NOT NULL,

			UNIQUE (dead_nation_id)
		);

		CREATE TABLE IF NOT EXISTS bookings (
				id UUID PRIMARY KEY,
				show_id UUID NOT NULL,
				number_of_tickets INT NOT NULL,
				customer_email VARCHAR(255) NOT NULL,
				FOREIGN KEY (show_id) REFERENCES shows(id)
		);
	`)
	if err != nil {
		return fmt.Errorf("failed to initialize schema: %w", err)
	}

	return nil
}

type DBTX interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	GetContext(ctx context.Context, dest any, query string, args ...any) error
	SelectContext(ctx context.Context, dest any, query string, args ...any) error
}

func (d *DB) WithTx(ctx context.Context) DBTX {
	sqlxTx, ok := ctx.Value(tx.ContextKey).(*sqlx.Tx)
	if !ok {
		return d.DB
	}
	return sqlxTx
}
