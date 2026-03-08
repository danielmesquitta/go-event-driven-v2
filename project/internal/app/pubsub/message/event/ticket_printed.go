package event

import (
	"context"
	"tickets/internal/app/pubsub/message"
)

type TicketPrinted struct {
	Header   message.Header `json:"header"`
	TicketID string         `json:"ticket_id"`
	FileName string         `json:"file_name"`
}

func NewTicketPrinted(
	ctx context.Context,
	ticketID string,
	fileName string,
) *TicketPrinted {
	return &TicketPrinted{
		Header:   message.NewHeader(ctx),
		TicketID: ticketID,
		FileName: fileName,
	}
}

func (e *TicketPrinted) GetHeader() message.Header {
	return e.Header
}

var _ Event = (*TicketPrinted)(nil)
