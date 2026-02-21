package event

import "tickets/internal/domain/entity"

type TicketBookingConfirmed struct {
	Header        EventHeader  `json:"header"`
	TicketID      string       `json:"ticket_id"`
	CustomerEmail string       `json:"customer_email"`
	Price         entity.Money `json:"price"`
}

func NewTicketBookingConfirmed(ticketID string, customerEmail string, price entity.Money) *TicketBookingConfirmed {
	return &TicketBookingConfirmed{
		Header:        NewEventHeader(),
		TicketID:      ticketID,
		CustomerEmail: customerEmail,
		Price:         price,
	}
}

func (e *TicketBookingConfirmed) GetHeader() EventHeader {
	return e.Header
}

func (e *TicketBookingConfirmed) SetDefaults() {
	if e.Header.ID == "" || e.Header.PublishedAt.IsZero() {
		e.Header = NewEventHeader()
	}

	if e.Price.Currency == "" {
		e.Price.Currency = "USD"
	}
}

var _ Event = (*TicketBookingConfirmed)(nil)
