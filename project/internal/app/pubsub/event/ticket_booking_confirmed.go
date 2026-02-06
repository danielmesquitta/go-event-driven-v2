package event

import "tickets/internal/domain/entity"

const TopicTicketBookingConfirmed Topic = "TicketBookingConfirmed"

type TicketBookingConfirmed struct {
	Header        EventHeader  `json:"header"`
	TicketID      string       `json:"ticket_id"`
	CustomerEmail string       `json:"customer_email"`
	Price         entity.Money `json:"price"`
}
