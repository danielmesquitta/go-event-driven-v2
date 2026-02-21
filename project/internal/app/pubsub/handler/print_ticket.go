package handler

import (
	"context"
	"tickets/internal/app/pubsub/event"
	"tickets/internal/domain/usecase"
)

type PrintTicket struct {
	printTicketUseCase *usecase.PrintTicket
}

func NewPrintTicket(printTicketUseCase *usecase.PrintTicket) *PrintTicket {
	return &PrintTicket{printTicketUseCase: printTicketUseCase}
}

func (p *PrintTicket) Handle(ctx context.Context, event *event.TicketBookingConfirmed) error {
	return p.printTicketUseCase.Execute(ctx, usecase.PrintTicketInput{
		TicketID: event.TicketID,
		Price:    event.Price,
	})
}
