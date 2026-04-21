package opsbooking

import (
	"context"
	"fmt"

	"tickets/internal/domain/entity"
	"tickets/internal/pkg/validator"
	"tickets/internal/provider/repo"
)

type OnTicketRefunded struct {
	opsBookingRepo repo.OpsBookingRepo
}

func NewOnTicketRefunded(opsBookingRepo repo.OpsBookingRepo) *OnTicketRefunded {
	return &OnTicketRefunded{opsBookingRepo: opsBookingRepo}
}

type OnTicketRefundedInput struct {
	TicketID string `json:"ticket_id" validate:"required"`
}

func (uc *OnTicketRefunded) Execute(ctx context.Context, in OnTicketRefundedInput) error {
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

			ticket.Status = "refunded"
			op.Tickets[in.TicketID] = ticket

			return op, nil
		},
	)
}
