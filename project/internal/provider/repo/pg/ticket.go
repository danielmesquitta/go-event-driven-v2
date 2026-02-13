package pg

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
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
	ID            string `db:"ticket_id"`
	PriceAmount   string `db:"price_amount"`
	PriceCurrency string `db:"price_currency"`
	CustomerEmail string `db:"customer_email"`
}

func (r *TicketRepo) Create(ctx context.Context, ticket *entity.Ticket) error {
	_, err := r.db.DB.ExecContext(ctx, `
		INSERT INTO tickets (ticket_id, price_amount, price_currency, customer_email)
		VALUES ($1, $2, $3, $4)
	`, ticket.ID, ticket.Price.Amount, ticket.Price.Currency, ticket.CustomerEmail)
	if err != nil {
		return fmt.Errorf("failed to create ticket: %w", err)
	}
	return nil
}

func (r *TicketRepo) Get(ctx context.Context, id string) (*entity.Ticket, error) {
	var ticket Ticket
	err := r.db.DB.GetContext(ctx, &ticket, `
		SELECT ticket_id, price_amount, price_currency, customer_email
		FROM tickets
		WHERE ticket_id = $1
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
	}, nil
}

func (r *TicketRepo) Delete(ctx context.Context, id string) error {
	_, err := r.db.DB.ExecContext(ctx, `
		DELETE FROM tickets
		WHERE ticket_id = $1
	`, id)
	if err != nil {
		return fmt.Errorf("failed to delete ticket: %w", err)
	}
	return nil
}

var _ repo.TicketRepo = &TicketRepo{}
