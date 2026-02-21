package handler

import (
	"context"
	"tickets/internal/app/pubsub/event"
	"tickets/internal/domain/usecase"
)

type DeleteTicket struct {
	deleteTicketUseCase *usecase.DeleteTicket
}

func NewDeleteTicket(deleteTicketUseCase *usecase.DeleteTicket) *DeleteTicket {
	return &DeleteTicket{deleteTicketUseCase: deleteTicketUseCase}
}

func (c *DeleteTicket) Handle(ctx context.Context, event *event.TicketBookingCanceled) error {
	err := c.deleteTicketUseCase.Execute(ctx, usecase.DeleteTicketInput{
		TicketID: event.TicketID,
	})
	if err != nil {
		return err
	}
	return nil
}
