package repo

import (
	"context"
	"time"

	"tickets/internal/domain/entity"
)

type CreateOpsBookingInput struct {
	BookingID string    `json:"booking_id" validate:"required"`
	BookedAt  time.Time `json:"booked_at"  validate:"required"`
}

type OpsBookingRepo interface {
	Create(ctx context.Context, in CreateOpsBookingInput) error
	GetByBookingID(ctx context.Context, bookingID string) (*entity.OpsBooking, error)
	GetByTicketID(ctx context.Context, ticketID string) (*entity.OpsBooking, error)

	// UpdateByBookingID atomically updates the read model identified by booking ID.
	// The update function receives the current state and returns the updated one.
	UpdateByBookingID(
		ctx context.Context,
		bookingID string,
		updateFn func(entity.OpsBooking) (entity.OpsBooking, error),
	) error

	// UpdateByTicketID atomically updates the read model that contains the given ticket.
	// The update function receives the full ops booking (not just the ticket) and returns it updated.
	UpdateByTicketID(
		ctx context.Context,
		ticketID string,
		updateFn func(entity.OpsBooking) (entity.OpsBooking, error),
	) error
}
