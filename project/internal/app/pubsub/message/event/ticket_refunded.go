package event

import (
	"context"
	"tickets/internal/app/pubsub/message"
)

type TicketRefunded struct {
	Header   message.Header `json:"header"`
	TicketID string         `json:"ticket_id"`
}

func NewTicketRefunded(
	ctx context.Context,
	ticketID string,
) *TicketRefunded {
	return &TicketRefunded{
		Header:   message.NewHeader(ctx),
		TicketID: ticketID,
	}
}

func (e *TicketRefunded) GetHeader() message.Header {
	return e.Header
}

var _ Event = (*TicketRefunded)(nil)
