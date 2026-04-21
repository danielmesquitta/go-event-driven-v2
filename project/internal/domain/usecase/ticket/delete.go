package ticket

import (
	"context"
	"tickets/internal/pkg/validator"
	"tickets/internal/provider/repo"
)

type Delete struct {
	ticketRepo repo.TicketRepo
}

func NewDelete(ticketRepo repo.TicketRepo) *Delete {
	return &Delete{ticketRepo: ticketRepo}
}

type DeleteInput struct {
	TicketID string `json:"ticket_id" validate:"required"`
}

func (c *Delete) Execute(ctx context.Context, in DeleteInput) error {
	err := validator.Validate(ctx, in)
	if err != nil {
		return err
	}

	// The repository performs a soft delete and returns an error when the ticket
	// does not exist yet. That way, if TicketBookingCanceled arrives before
	// TicketBookingConfirmed, the message is nacked and redelivered later.
	// For redeliveries of an already soft-deleted ticket, the UPDATE still
	// matches the row, so the operation stays idempotent.
	return c.ticketRepo.Delete(ctx, in.TicketID)
}
