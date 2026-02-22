package event

import (
	"context"
	"tickets/internal/domain/entity"
)

type TicketBookingConfirmed struct {
	Header        EventHeader  `json:"header"`
	TicketID      string       `json:"ticket_id"`
	CustomerEmail string       `json:"customer_email"`
	Price         entity.Money `json:"price"`
}

func NewTicketBookingConfirmed(
	ctx context.Context,
	ticketID string,
	customerEmail string,
	price entity.Money,
) *TicketBookingConfirmed {
	return &TicketBookingConfirmed{
		Header:        NewEventHeader(ctx),
		TicketID:      ticketID,
		CustomerEmail: customerEmail,
		Price:         price,
	}
}

func (e *TicketBookingConfirmed) GetHeader() EventHeader {
	return e.Header
}

var _ Event = (*TicketBookingConfirmed)(nil)
