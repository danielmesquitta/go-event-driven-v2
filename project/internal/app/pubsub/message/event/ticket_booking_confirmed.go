package event

import (
	"context"
	"tickets/internal/app/pubsub/message"
	"tickets/internal/domain/entity"
)

type TicketBookingConfirmed struct {
	Header        message.Header `json:"header"`
	TicketID      string         `json:"ticket_id"`
	CustomerEmail string         `json:"customer_email"`
	Price         entity.Money   `json:"price"`
}

func NewTicketBookingConfirmed(
	ctx context.Context,
	ticketID string,
	customerEmail string,
	price entity.Money,
) *TicketBookingConfirmed {
	return &TicketBookingConfirmed{
		Header:        message.NewHeader(ctx),
		TicketID:      ticketID,
		CustomerEmail: customerEmail,
		Price:         price,
	}
}

func (e *TicketBookingConfirmed) GetHeader() message.Header {
	return e.Header
}

var _ Event = (*TicketBookingConfirmed)(nil)
