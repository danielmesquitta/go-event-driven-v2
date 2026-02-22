package usecase

import (
	"context"
	"tickets/internal/pkg/validator"
	"tickets/internal/provider/repo"
)

type DeleteTicket struct {
	ticketRepo repo.TicketRepo
}

func NewDeleteTicket(ticketRepo repo.TicketRepo) *DeleteTicket {
	return &DeleteTicket{ticketRepo: ticketRepo}
}

type DeleteTicketInput struct {
	TicketID string `json:"ticket_id" validate:"required"`
}

func (c *DeleteTicket) Execute(ctx context.Context, in DeleteTicketInput) error {
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
