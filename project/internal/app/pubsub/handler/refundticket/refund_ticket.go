package refundticket

import (
	"context"
	"tickets/internal/app/pubsub/message/cmd"
	"tickets/internal/domain/usecase/ticket"
)

type RefundTicket struct {
	refundTicketUseCase *ticket.Refund
}

func NewRefundTicket(refundTicketUseCase *ticket.Refund) *RefundTicket {
	return &RefundTicket{refundTicketUseCase: refundTicketUseCase}
}

func (c *RefundTicket) Handle(ctx context.Context, cmd *cmd.RefundTicket) error {
	err := c.refundTicketUseCase.Execute(ctx, ticket.RefundInput{
		TicketID: cmd.TicketID,
	})
	if err != nil {
		return err
	}
	return nil
}
