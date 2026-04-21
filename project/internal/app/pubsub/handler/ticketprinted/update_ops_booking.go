package ticketprinted

import (
	"context"

	"tickets/internal/app/pubsub/message/event"
	"tickets/internal/domain/usecase/opsbooking"
)

type UpdateOpsBooking struct {
	useCase *opsbooking.OnTicketPrinted
}

func NewUpdateOpsBooking(useCase *opsbooking.OnTicketPrinted) *UpdateOpsBooking {
	return &UpdateOpsBooking{useCase: useCase}
}

func (h *UpdateOpsBooking) Handle(ctx context.Context, e *event.TicketPrinted) error {
	return h.useCase.Execute(ctx, opsbooking.OnTicketPrintedInput{
		TicketID:  e.TicketID,
		FileName:  e.FileName,
		PrintedAt: e.Header.PublishedAt,
	})
}
