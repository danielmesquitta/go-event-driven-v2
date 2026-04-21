package opsbooking

import (
	"context"
	"fmt"
	"tickets/internal/domain/errs"
	"tickets/internal/pkg/validator"
	"tickets/internal/provider/repo"
	"time"
)

type Update struct {
	opsBookingRepo repo.OpsBookingRepo
}

func NewUpdate(opsBookingRepo repo.OpsBookingRepo) *Update {
	return &Update{opsBookingRepo: opsBookingRepo}
}

type UpdateInput struct {
	BookingID string    `json:"booking_id" validate:"required"`
	BookedAt  time.Time `json:"booked_at" validate:"required"`
}

func (uc *Update) Execute(ctx context.Context, in UpdateInput) error {
	err := validator.Validate(ctx, in)
	if err != nil {
		return fmt.Errorf("failed to validate update ops booking input: %w", err)
	}

	opsBooking, err := uc.opsBookingRepo.GetByBookingID(ctx, in.BookingID)
	if err != nil {
		return fmt.Errorf("failed to get ops booking: %w", err)
	}
	if opsBooking == nil {
		return errs.ErrNotFound.New(
			errs.WithMessage("ops booking not found"),
			errs.WithMetadata(errs.MetadataDataKey, map[string]string{
				"booking_id": in.BookingID,
			}),
		)
	}

	return nil
}
