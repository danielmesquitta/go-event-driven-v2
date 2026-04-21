package bookingconfirmed

import (
	"context"
	"tickets/internal/app/pubsub/message/event"
	"tickets/internal/domain/usecase/ticket"
)

type UpdateOpsBooking struct {
	printTicketUseCase *ticket.Print
}

func NewUpdateOpsBooking(printTicketUseCase *ticket.Print) *UpdateOpsBooking {
	return &UpdateOpsBooking{printTicketUseCase: printTicketUseCase}
}

func (p *UpdateOpsBooking) Handle(ctx context.Context, event *event.TicketBookingConfirmed) error {
	return p.printTicketUseCase.Execute(ctx, ticket.PrintInput{
		TicketID: event.TicketID,
		Price:    event.Price,
	})
}
