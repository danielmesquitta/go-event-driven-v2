package handler

import (
	"context"
	"tickets/internal/app/pubsub/event"
	"tickets/internal/provider/repo"
)

type DeleteTicket struct {
	ticketRepo repo.TicketRepo
}

func NewDeleteTicket(ticketRepo repo.TicketRepo) *DeleteTicket {
	return &DeleteTicket{ticketRepo: ticketRepo}
}

func (c *DeleteTicket) Handle(ctx context.Context, event *event.TicketBookingCanceled) error {
	ticket, err := c.ticketRepo.Get(ctx, event.TicketID)
	if err != nil {
		return err
	}
	if ticket == nil {
		return nil
	}

	err = c.ticketRepo.Delete(ctx, event.TicketID)
	if err != nil {
		return err
	}

	return nil
}
