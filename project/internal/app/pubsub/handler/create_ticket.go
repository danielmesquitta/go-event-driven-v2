package handler

import (
	"context"
	"tickets/internal/app/pubsub/event"
	"tickets/internal/domain/entity"
	"tickets/internal/provider/repo"
)

type CreateTicket struct {
	ticketRepo repo.TicketRepo
}

func NewCreateTicket(ticketRepo repo.TicketRepo) *CreateTicket {
	return &CreateTicket{ticketRepo: ticketRepo}
}

func (c *CreateTicket) Handle(ctx context.Context, event *event.TicketBookingConfirmed) error {
	err := c.ticketRepo.Create(ctx, &entity.Ticket{
		ID:            event.TicketID,
		Price:         event.Price,
		CustomerEmail: event.CustomerEmail,
	})
	if err != nil {
		return err
	}
	return nil
}
