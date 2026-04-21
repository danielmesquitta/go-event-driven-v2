package opsbooking

import (
	"context"
	"fmt"
	"time"

	"tickets/internal/domain/entity"
	"tickets/internal/provider/repo"
)

type List struct {
	opsBookingRepo repo.OpsBookingRepo
}

func NewList(opsBookingRepo repo.OpsBookingRepo) *List {
	return &List{opsBookingRepo: opsBookingRepo}
}

type ListInput struct {
	// ReceiptIssueDate, if set, filters results to bookings that have at least
	// one ticket with a receipt issued on this date.
	ReceiptIssueDate *time.Time
}

func (uc *List) Execute(ctx context.Context, in ListInput) ([]entity.OpsBooking, error) {
	bookings, err := uc.opsBookingRepo.List(ctx, repo.ListOpsBookingsFilter{
		ReceiptIssueDate: in.ReceiptIssueDate,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list ops bookings: %w", err)
	}

	return bookings, nil
}
