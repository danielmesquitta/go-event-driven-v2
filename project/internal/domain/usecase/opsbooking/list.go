package opsbooking

import (
	"context"
	"fmt"

	"tickets/internal/domain/entity"
	"tickets/internal/provider/repo"
)

type List struct {
	opsBookingRepo repo.OpsBookingRepo
}

func NewList(opsBookingRepo repo.OpsBookingRepo) *List {
	return &List{opsBookingRepo: opsBookingRepo}
}

func (uc *List) Execute(ctx context.Context) ([]entity.OpsBooking, error) {
	bookings, err := uc.opsBookingRepo.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list ops bookings: %w", err)
	}

	return bookings, nil
}
