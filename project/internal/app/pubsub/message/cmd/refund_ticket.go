package cmd

import (
	"context"
	"tickets/internal/app/pubsub/message"
)

type RefundTicket struct {
	Header   message.Header `json:"header"`
	TicketID string         `json:"ticket_id"`
}

func NewRefundTicket(
	ctx context.Context,
	ticketID string,
) *RefundTicket {
	return &RefundTicket{
		Header:   message.NewHeader(ctx),
		TicketID: ticketID,
	}
}

func (e *RefundTicket) GetHeader() message.Header {
	return e.Header
}

var _ Command = (*RefundTicket)(nil)
