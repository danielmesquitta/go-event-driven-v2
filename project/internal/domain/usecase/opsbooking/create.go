package opsbooking

import (
	"context"
	"fmt"
	"tickets/internal/pkg/tx"
	"tickets/internal/pkg/validator"
	"tickets/internal/provider/repo"
	"time"
)

type Create struct {
	tx             tx.Transaction
	opsBookingRepo repo.OpsBookingRepo
}

func NewCreate(
	tx tx.Transaction,
	opsBookingRepo repo.OpsBookingRepo,
) *Create {
	return &Create{
		tx:             tx,
		opsBookingRepo: opsBookingRepo,
	}
}

type CreateInput struct {
	BookingID string    `json:"booking_id" validate:"required"`
	BookedAt  time.Time `json:"booked_at" validate:"required"`
}

func (uc *Create) Execute(ctx context.Context, in CreateInput) error {
	err := validator.Validate(ctx, in)
	if err != nil {
		return fmt.Errorf("failed to validate create ops booking input: %w", err)
	}

	err = uc.opsBookingRepo.Create(ctx, repo.CreateOpsBookingInput{
		BookingID: in.BookingID,
		BookedAt:  in.BookedAt,
	})
	if err != nil {
		return fmt.Errorf("failed to create ops booking: %w", err)
	}

	return nil
}
