package opsbooking

import (
	"context"
	"fmt"

	"tickets/internal/domain/entity"
	"tickets/internal/domain/errs"
	"tickets/internal/pkg/validator"
	"tickets/internal/provider/repo"
)

type Get struct {
	opsBookingRepo repo.OpsBookingRepo
}

func NewGet(opsBookingRepo repo.OpsBookingRepo) *Get {
	return &Get{opsBookingRepo: opsBookingRepo}
}

type GetInput struct {
	BookingID string `json:"booking_id" validate:"required"`
}

func (uc *Get) Execute(ctx context.Context, in GetInput) (entity.OpsBooking, error) {
	if err := validator.Validate(ctx, in); err != nil {
		return entity.OpsBooking{}, fmt.Errorf("failed to validate get ops booking input: %w", err)
	}

	booking, err := uc.opsBookingRepo.GetByBookingID(ctx, in.BookingID)
	if err != nil {
		return entity.OpsBooking{}, fmt.Errorf("failed to get ops booking: %w", err)
	}
	if booking == nil {
		return entity.OpsBooking{}, errs.ErrNotFound.New()
	}

	return *booking, nil
}
