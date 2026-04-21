package repo

import (
	"context"
	"tickets/internal/domain/entity"
	"time"
)

type CreateOpsBookingInput struct {
	BookingID string    `json:"booking_id" validate:"required"`
	BookedAt  time.Time `json:"booked_at"  validate:"required"`
}

type UpdateOpsTicketInput struct {
}

type OpsBookingRepo interface {
	Create(ctx context.Context, in CreateOpsBookingInput) error
	Update(ctx context.Context, op entity.OpsBooking) error
	GetByBookingID(ctx context.Context, bookingID string) (*entity.OpsBooking, error)
	GetByTicketID(ctx context.Context, ticketID string) (*entity.OpsBooking, error)
}
