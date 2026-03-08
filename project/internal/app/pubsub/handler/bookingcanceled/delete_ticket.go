package bookingcanceled

import (
	"context"
	"tickets/internal/app/pubsub/message/event"
	"tickets/internal/domain/usecase/ticket"
)

type DeleteTicket struct {
	deleteTicketUseCase *ticket.Delete
}

func NewDeleteTicket(deleteTicketUseCase *ticket.Delete) *DeleteTicket {
	return &DeleteTicket{deleteTicketUseCase: deleteTicketUseCase}
}

func (c *DeleteTicket) Handle(ctx context.Context, event *event.TicketBookingCanceled) error {
	err := c.deleteTicketUseCase.Execute(ctx, ticket.DeleteInput{
		TicketID: event.TicketID,
	})
	if err != nil {
		return err
	}
	return nil
}
