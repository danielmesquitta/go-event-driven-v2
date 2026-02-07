package event

import "tickets/internal/domain/entity"

const TopicTicketBookingConfirmed Topic = "TicketBookingConfirmed"

type TicketBookingConfirmed struct {
	Header        EventHeader  `json:"header"`
	TicketID      string       `json:"ticket_id"`
	CustomerEmail string       `json:"customer_email"`
	Price         entity.Money `json:"price"`
}

func (e *TicketBookingConfirmed) GetHeader() EventHeader {
	return e.Header
}

// SetDefaults sets the default values for the event, in case they are not set.
func (e *TicketBookingConfirmed) SetDefaults() {
	if e.Header.ID == "" || e.Header.PublishedAt.IsZero() {
		e.Header = NewEventHeader()
	}

	if e.Price.Currency == "" {
		e.Price.Currency = "USD"
	}
}

var _ Event = (*TicketBookingConfirmed)(nil)
