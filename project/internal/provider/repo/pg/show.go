package pg

import (
	"context"
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
