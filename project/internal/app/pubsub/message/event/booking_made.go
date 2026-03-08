package event

import (
	"context"
	"tickets/internal/app/pubsub/message"
)

type BookingMade struct {
	Header          message.Header `json:"header"`
	BookingID       string         `json:"booking_id"`
	NumberOfTickets int            `json:"number_of_tickets"`
	CustomerEmail   string         `json:"customer_email"`
	ShowID          string         `json:"show_id"`
}

func NewBookingMade(
	ctx context.Context,
	bookingID string,
	numberOfTickets int,
	customerEmail string,
	showID string,
) *BookingMade {
	return &BookingMade{
		Header:          message.NewHeader(ctx),
		BookingID:       bookingID,
		NumberOfTickets: numberOfTickets,
		CustomerEmail:   customerEmail,
		ShowID:          showID,
	}
}

func (e *BookingMade) GetHeader() message.Header {
	return e.Header
}

var _ Event = (*BookingMade)(nil)
