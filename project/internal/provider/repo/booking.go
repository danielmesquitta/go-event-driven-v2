package repo

import (
	"context"
	"tickets/internal/domain/entity"
)

type BookingRepo interface {
	Create(ctx context.Context, booking *entity.Booking) error
}
