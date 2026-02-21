package event

import "tickets/internal/domain/entity"

type TicketBookingCanceled struct {
	Header        EventHeader  `json:"header"`
	TicketID      string       `json:"ticket_id"`
	CustomerEmail string       `json:"customer_email"`
	Price         entity.Money `json:"price"`
}

func NewTicketBookingCanceled(ticketID string, customerEmail string, price entity.Money) *TicketBookingCanceled {
	return &TicketBookingCanceled{
		Header:        NewEventHeader(),
		TicketID:      ticketID,
		CustomerEmail: customerEmail,
		Price:         price,
	}
}

func (e *TicketBookingCanceled) GetHeader() EventHeader {
	return e.Header
}

func (e *TicketBookingCanceled) SetDefaults() {
	if e.Header.ID == "" || e.Header.PublishedAt.IsZero() {
		e.Header = NewEventHeader()
	}

	if e.Price.Currency == "" {
		e.Price.Currency = "USD"
	}
}

var _ Event = (*TicketBookingCanceled)(nil)
