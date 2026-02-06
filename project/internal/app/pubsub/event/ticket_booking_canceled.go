package event

import "tickets/internal/domain/entity"

const TopicTicketBookingCanceled Topic = "TicketBookingCanceled"

type TicketBookingCanceled struct {
	Header        EventHeader  `json:"header"`
	TicketID      string       `json:"ticket_id"`
	CustomerEmail string       `json:"customer_email"`
	Price         entity.Money `json:"price"`
}

func (e *TicketBookingCanceled) GetHeader() EventHeader {
	return e.Header
}

var _ Event = (*TicketBookingCanceled)(nil)
