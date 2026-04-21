package ticketreceiptissued

import (
	"context"

	"tickets/internal/app/pubsub/message/event"
	"tickets/internal/domain/usecase/opsbooking"
)

type UpdateOpsBooking struct {
	useCase *opsbooking.OnTicketReceiptIssued
}

func NewUpdateOpsBooking(useCase *opsbooking.OnTicketReceiptIssued) *UpdateOpsBooking {
	return &UpdateOpsBooking{useCase: useCase}
}

func (h *UpdateOpsBooking) Handle(ctx context.Context, e *event.TicketReceiptIssued) error {
	return h.useCase.Execute(ctx, opsbooking.OnTicketReceiptIssuedInput{
		TicketID:      e.TicketID,
		ReceiptNumber: e.ReceiptNumber,
		IssuedAt:      e.IssuedAt,
	})
}
