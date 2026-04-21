package pg

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"tickets/internal/domain/entity"
	"tickets/internal/provider/db"
	"tickets/internal/provider/repo"
)

type TicketRepo struct {
	db *db.DB
}

func NewTicketRepo(db *db.DB) *TicketRepo {
	return &TicketRepo{db: db}
}

type Ticket struct {
	ID            string `db:"id"`
	PriceAmount   string `db:"price_amount"`
	PriceCurrency string `db:"price_currency"`
	CustomerEmail string `db:"customer_email"`
	BookingID     string `db:"booking_id"`
}

func (r *TicketRepo) Create(ctx context.Context, ticket *entity.Ticket) error {
	_, err := r.db.WithTx(ctx).ExecContext(ctx, `
		INSERT INTO tickets (id, price_amount, price_currency, customer_email, booking_id)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT DO NOTHING
	`, ticket.ID, ticket.Price.Amount, ticket.Price.Currency, ticket.CustomerEmail, ticket.BookingID)
	if err != nil {
		return fmt.Errorf("failed to create ticket: %w", err)
	}
	return nil
}

func (r *TicketRepo) Get(ctx context.Context, id string) (*entity.Ticket, error) {
	var ticket Ticket
	err := r.db.WithTx(ctx).GetContext(ctx, &ticket, `
		SELECT *
		FROM tickets
		WHERE id = $1
		LIMIT 1
	`, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get ticket: %w", err)
	}

	return &entity.Ticket{
		ID: ticket.ID,
		Price: entity.Money{
			Amount:   ticket.PriceAmount,
			Currency: ticket.PriceCurrency,
		},
		CustomerEmail: ticket.CustomerEmail,
		BookingID:     ticket.BookingID,
	}, nil
}

func (r *TicketRepo) List(ctx context.Context) ([]entity.Ticket, error) {
	var tickets []Ticket
	err := r.db.WithTx(ctx).SelectContext(ctx, &tickets, `
		SELECT *
		FROM tickets
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to list tickets: %w", err)
	}

	entities := make([]entity.Ticket, len(tickets))
	for i, ticket := range tickets {
		entities[i] = entity.Ticket{
			ID: ticket.ID,
			Price: entity.Money{
				Amount:   ticket.PriceAmount,
				Currency: ticket.PriceCurrency,
			},
			CustomerEmail: ticket.CustomerEmail,
			BookingID:     ticket.BookingID,
		}
	}

	slog.InfoContext(ctx, "tickets", "tickets", entities)

	return entities, nil
}

func (r *TicketRepo) Delete(ctx context.Context, id string) error {
	_, err := r.db.WithTx(ctx).ExecContext(ctx, `
		DELETE FROM tickets
		WHERE id = $1
	`, id)
	if err != nil {
		return fmt.Errorf("failed to delete ticket: %w", err)
	}
	return nil
}

var _ repo.TicketRepo = &TicketRepo{}
