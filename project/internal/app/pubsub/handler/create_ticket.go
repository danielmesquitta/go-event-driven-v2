package handler

import (
	"context"
	"tickets/internal/app/pubsub/event"
	"tickets/internal/domain/usecase"
)

type CreateTicket struct {
	createTicketUseCase *usecase.CreateTicket
}

func NewCreateTicket(createTicketUseCase *usecase.CreateTicket) *CreateTicket {
	return &CreateTicket{createTicketUseCase: createTicketUseCase}
}

func (c *CreateTicket) Handle(ctx context.Context, event *event.TicketBookingConfirmed) error {
	err := c.createTicketUseCase.Execute(ctx, usecase.CreateTicketInput{
		TicketID:      event.TicketID,
		Price:         event.Price,
		CustomerEmail: event.CustomerEmail,
	})
	if err != nil {
		return err
	}
	return nil
}
