package bookingconfirmed

import (
	"context"
	"tickets/internal/app/pubsub/message/event"
	"tickets/internal/domain/usecase/ticket"
)

type PrintTicket struct {
	printTicketUseCase *ticket.Print
}

func NewPrintTicket(printTicketUseCase *ticket.Print) *PrintTicket {
	return &PrintTicket{printTicketUseCase: printTicketUseCase}
}

func (p *PrintTicket) Handle(ctx context.Context, event *event.TicketBookingConfirmed) error {
	return p.printTicketUseCase.Execute(ctx, ticket.PrintInput{
		TicketID: event.TicketID,
		Price:    event.Price,
	})
}
