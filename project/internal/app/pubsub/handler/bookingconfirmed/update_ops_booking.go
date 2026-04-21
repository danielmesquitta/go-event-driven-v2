package bookingconfirmed

import (
	"context"

	"tickets/internal/app/pubsub/message/event"
	"tickets/internal/domain/usecase/opsbooking"
)

type UpdateOpsBooking struct {
	useCase *opsbooking.OnTicketBookingConfirmed
}

func NewUpdateOpsBooking(useCase *opsbooking.OnTicketBookingConfirmed) *UpdateOpsBooking {
	return &UpdateOpsBooking{useCase: useCase}
}

func (h *UpdateOpsBooking) Handle(ctx context.Context, e *event.TicketBookingConfirmed) error {
	return h.useCase.Execute(ctx, opsbooking.OnTicketBookingConfirmedInput{
		BookingID:     e.BookingID,
		TicketID:      e.TicketID,
		CustomerEmail: e.CustomerEmail,
		Price:         e.Price,
	})
}
