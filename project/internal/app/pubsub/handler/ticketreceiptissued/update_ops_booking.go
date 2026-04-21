package ticketreceiptissued

import (
	"context"
	"tickets/internal/app/pubsub/message/event"
)

type UpdateOpsBooking struct {
}

func NewUpdateOpsBooking() *UpdateOpsBooking {
	return &UpdateOpsBooking{}
}

func (c *UpdateOpsBooking) Handle(ctx context.Context, e *event.TicketRefunded) error {
	// err := c.refundTicketUseCase.Execute(ctx, ticket.RefundInput{
	// 	TicketID: cmd.TicketID,
	// })
	// if err != nil {
	// 	return err
	// }
	return nil
}
