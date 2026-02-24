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

	ticket, err := c.ticketRepo.Get(ctx, in.TicketID)
	if err != nil {
		return err
	}
	if ticket == nil {
		return nil
	}

	err = c.ticketRepo.Delete(ctx, in.TicketID)
	if err != nil {
		return err
	}

	return nil
}
