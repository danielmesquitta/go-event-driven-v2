package ticketrefunded

import (
	"context"

	"tickets/internal/app/pubsub/message/event"
	"tickets/internal/domain/usecase/opsbooking"
)

type UpdateOpsBooking struct {
	useCase *opsbooking.OnTicketRefunded
}

func NewUpdateOpsBooking(useCase *opsbooking.OnTicketRefunded) *UpdateOpsBooking {
	return &UpdateOpsBooking{useCase: useCase}
}

func (h *UpdateOpsBooking) Handle(ctx context.Context, e *event.TicketRefunded) error {
	return h.useCase.Execute(ctx, opsbooking.OnTicketRefundedInput{
		TicketID:   e.TicketID,
		RefundedAt: e.Header.PublishedAt,
	})
}
