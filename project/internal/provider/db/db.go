package db

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var timeout = 10 * time.Second

type DB struct {
	*sqlx.DB
}

func NewDB() *DB {
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

	return &DB{
		DB: db,
	}
}

func (d *DB) InitializeSchema(ctx context.Context) error {
	_, err := d.DB.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS tickets (
			ticket_id UUID PRIMARY KEY,
			price_amount   NUMERIC(10, 2) NOT NULL,
			price_currency CHAR(3) NOT NULL,
			customer_email VARCHAR(255) NOT NULL
		);
	`)
	if err != nil {
		return fmt.Errorf("failed to initialize schema: %w", err)
	}

	return nil
}
