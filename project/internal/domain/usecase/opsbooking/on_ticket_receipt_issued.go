package opsbooking

import (
	"context"
	"fmt"
	"time"

	"tickets/internal/domain/entity"
	"tickets/internal/pkg/validator"
	"tickets/internal/provider/repo"
)

type OnTicketReceiptIssued struct {
	opsBookingRepo repo.OpsBookingRepo
}

func NewOnTicketReceiptIssued(opsBookingRepo repo.OpsBookingRepo) *OnTicketReceiptIssued {
	return &OnTicketReceiptIssued{opsBookingRepo: opsBookingRepo}
}

type OnTicketReceiptIssuedInput struct {
	TicketID      string    `json:"ticket_id" validate:"required"`
	ReceiptNumber string    `json:"receipt_number" validate:"required"`
	IssuedAt      time.Time `json:"issued_at" validate:"required"`
}

func (uc *OnTicketReceiptIssued) Execute(ctx context.Context, in OnTicketReceiptIssuedInput) error {
	if err := validator.Validate(ctx, in); err != nil {
		return fmt.Errorf("failed to validate input: %w", err)
	}

	return uc.opsBookingRepo.UpdateByTicketID(
		ctx,
		in.TicketID,
		func(op entity.OpsBooking) (entity.OpsBooking, error) {
			ticket, ok := op.Tickets[in.TicketID]
			if !ok {
				return op, fmt.Errorf("ticket %s not found in ops booking", in.TicketID)
			}

			ticket.ReceiptIssuedAt = in.IssuedAt
			ticket.ReceiptNumber = in.ReceiptNumber
			op.Tickets[in.TicketID] = ticket

			return op, nil
		},
	)
}
