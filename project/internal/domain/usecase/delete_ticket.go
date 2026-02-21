package usecase

import (
	"context"
	"tickets/internal/provider/repo"
)

type DeleteTicket struct {
	ticketRepo repo.TicketRepo
}

func NewDeleteTicket(ticketRepo repo.TicketRepo) *DeleteTicket {
	return &DeleteTicket{ticketRepo: ticketRepo}
}

type DeleteTicketInput struct {
	TicketID string
}

func (c *DeleteTicket) Execute(ctx context.Context, in DeleteTicketInput) error {
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
