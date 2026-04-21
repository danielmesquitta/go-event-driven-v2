package opsbooking

import (
	"context"
	"fmt"

	"tickets/internal/domain/entity"
	"tickets/internal/pkg/log"
	"tickets/internal/pkg/validator"
	"tickets/internal/provider/repo"
)

type OnTicketBookingConfirmed struct {
	opsBookingRepo repo.OpsBookingRepo
}

func NewOnTicketBookingConfirmed(opsBookingRepo repo.OpsBookingRepo) *OnTicketBookingConfirmed {
	return &OnTicketBookingConfirmed{opsBookingRepo: opsBookingRepo}
}

type OnTicketBookingConfirmedInput struct {
	BookingID     string       `json:"booking_id" validate:"required"`
	TicketID      string       `json:"ticket_id" validate:"required"`
	CustomerEmail string       `json:"customer_email" validate:"required"`
	Price         entity.Money `json:"price" validate:"required"`
}

func (uc *OnTicketBookingConfirmed) Execute(ctx context.Context, in OnTicketBookingConfirmedInput) error {
	if err := validator.Validate(ctx, in); err != nil {
		return fmt.Errorf("failed to validate input: %w", err)
	}

	return uc.opsBookingRepo.UpdateByBookingID(
		ctx,
		in.BookingID,
		func(op entity.OpsBooking) (entity.OpsBooking, error) {
			ticket, ok := op.Tickets[in.TicketID]
			if !ok {
				log.New(ctx).With("ticket_id", in.TicketID).Debug("Creating ticket in ops booking read model")
			}

			ticket.PriceAmount = in.Price.Amount
			ticket.PriceCurrency = in.Price.Currency
			ticket.CustomerEmail = in.CustomerEmail
			ticket.Status = "confirmed"
			op.Tickets[in.TicketID] = ticket

			return op, nil
		},
	)
}
