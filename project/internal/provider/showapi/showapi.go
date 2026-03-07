package showapi

import "context"

type Booking struct {
	BookingID       string `json:"booking_id"`
	EventID         string `json:"event_id"`
	NumberOfTickets int    `json:"number_of_tickets"`
	CustomerEmail   string `json:"customer_email"`
}

type ShowAPI interface {
	PostTicketBooking(ctx context.Context, booking Booking) error
}
