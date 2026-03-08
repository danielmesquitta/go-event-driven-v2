package pg

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"tickets/internal/domain/entity"
	"tickets/internal/provider/db"
)

type ShowRepo struct {
	db *db.DB
}

func NewShowRepo(db *db.DB) *ShowRepo {
	return &ShowRepo{db: db}
}

func (r *ShowRepo) Create(ctx context.Context, show *entity.Show) error {
	_, err := r.db.WithTx(ctx).ExecContext(ctx, `
		INSERT INTO shows (id, dead_nation_id, number_of_tickets, start_time, title, venue)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT DO NOTHING
	`, show.ID, show.DeadNationID, show.NumberOfTickets, show.StartTime, show.Title, show.Venue)
	if err != nil {
		return fmt.Errorf("failed to create show: %w", err)
	}
	return nil
}

func (r *ShowRepo) Get(ctx context.Context, id string) (*entity.Show, error) {
	var show entity.Show
	err := r.db.WithTx(ctx).GetContext(ctx, &show, `
		SELECT *
		FROM shows
		WHERE id = $1
		LIMIT 1
	`, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get show: %w", err)
	}
	return &show, nil
}

func (r *ShowRepo) AvailableTickets(ctx context.Context, id string) (int, error) {
	var availableTickets int
	err := r.db.WithTx(ctx).GetContext(ctx, &availableTickets, `
		SELECT shows.number_of_tickets - COALESCE(SUM(bookings.number_of_tickets), 0)
		FROM shows
		LEFT JOIN bookings ON shows.id = bookings.show_id
		WHERE shows.id = $1
		GROUP BY shows.id
	`, id)
	if err != nil {
		return 0, fmt.Errorf("failed to get available tickets: %w", err)
	}
	return availableTickets, nil
}
