package opsbooking

import (
	"context"
	"fmt"
	"time"

	"tickets/internal/domain/entity"
	"tickets/internal/pkg/validator"
	"tickets/internal/provider/repo"
)

type OnTicketPrinted struct {
	opsBookingRepo repo.OpsBookingRepo
}

func NewOnTicketPrinted(opsBookingRepo repo.OpsBookingRepo) *OnTicketPrinted {
	return &OnTicketPrinted{opsBookingRepo: opsBookingRepo}
}

type OnTicketPrintedInput struct {
	TicketID  string    `json:"ticket_id" validate:"required"`
	FileName  string    `json:"file_name" validate:"required"`
	PrintedAt time.Time `json:"printed_at" validate:"required"`
}

func (uc *OnTicketPrinted) Execute(ctx context.Context, in OnTicketPrintedInput) error {
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

			ticket.PrintedAt = in.PrintedAt
			ticket.PrintedFileName = in.FileName
			op.Tickets[in.TicketID] = ticket

			return op, nil
		},
	)
}
