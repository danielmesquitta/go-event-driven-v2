package bookingconfirmed

import (
	"context"
	"tickets/internal/app/pubsub/message/event"
	"tickets/internal/domain/usecase/ticket"
)

type CreateTicket struct {
	createTicketUseCase *ticket.Create
}

func NewCreateTicket(createTicketUseCase *ticket.Create) *CreateTicket {
	return &CreateTicket{createTicketUseCase: createTicketUseCase}
}

func (c *CreateTicket) Handle(ctx context.Context, event *event.TicketBookingConfirmed) error {
	err := c.createTicketUseCase.Execute(ctx, ticket.CreateInput{
		TicketID:      event.TicketID,
		Price:         event.Price,
		CustomerEmail: event.CustomerEmail,
	})
	if err != nil {
		return err
	}
	return nil
}
