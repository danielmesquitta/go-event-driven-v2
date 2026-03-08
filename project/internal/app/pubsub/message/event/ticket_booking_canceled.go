package event

import (
	"context"
	"tickets/internal/app/pubsub/message"
	"tickets/internal/domain/entity"
)

type TicketBookingCanceled struct {
	Header        message.Header `json:"header"`
	TicketID      string         `json:"ticket_id"`
	CustomerEmail string         `json:"customer_email"`
	Price         entity.Money   `json:"price"`
}

func NewTicketBookingCanceled(
	ctx context.Context,
	ticketID string,
	customerEmail string,
	price entity.Money,
) *TicketBookingCanceled {
	return &TicketBookingCanceled{
		Header:        message.NewHeader(ctx),
		TicketID:      ticketID,
		CustomerEmail: customerEmail,
		Price:         price,
	}
}

func (e *TicketBookingCanceled) GetHeader() message.Header {
	return e.Header
}

var _ Event = (*TicketBookingCanceled)(nil)
