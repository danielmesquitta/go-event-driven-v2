package pg

import (
	"context"
	"fmt"
	"tickets/internal/domain/entity"
	"tickets/internal/provider/db"
)

type BookingRepo struct {
	db *db.DB
}

func NewBookingRepo(db *db.DB) *BookingRepo {
	return &BookingRepo{db: db}
}

func (r *BookingRepo) Create(ctx context.Context, booking *entity.Booking) error {
	_, err := r.db.DB.ExecContext(ctx, `
		INSERT INTO bookings (id, show_id, number_of_tickets, customer_email)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT DO NOTHING
	`, booking.ID, booking.ShowID, booking.NumberOfTickets, booking.CustomerEmail)
	if err != nil {
		return fmt.Errorf("failed to create booking: %w", err)
	}
	return nil
}
