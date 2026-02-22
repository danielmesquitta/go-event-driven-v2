package event

import "context"

type TicketPrinted struct {
	Header   EventHeader `json:"header"`
	TicketID string      `json:"ticket_id"`
	FileName string      `json:"file_name"`
}

func NewTicketPrinted(
	ctx context.Context,
	ticketID string,
	fileName string,
) *TicketPrinted {
	return &TicketPrinted{
		Header:   NewEventHeader(ctx),
		TicketID: ticketID,
		FileName: fileName,
	}
}

func (e *TicketPrinted) GetHeader() EventHeader {
	return e.Header
}

var _ Event = (*TicketPrinted)(nil)
