package event

type TicketPrinted struct {
	Header   EventHeader `json:"header"`
	TicketID string      `json:"ticket_id"`
	FileName string      `json:"file_name"`
}

func NewTicketPrinted(ticketID string, fileName string) *TicketPrinted {
	return &TicketPrinted{
		Header:   NewEventHeader(),
		TicketID: ticketID,
		FileName: fileName,
	}
}

func (e *TicketPrinted) GetHeader() EventHeader {
	return e.Header
}

var _ Event = (*TicketPrinted)(nil)
