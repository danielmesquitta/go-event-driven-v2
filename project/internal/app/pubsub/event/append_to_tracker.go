package event

import "tickets/internal/domain/entity"

const TopicAppendToTracker Topic = "append-to-tracker"

type AppendToTrackerEvent struct {
	TicketID      string       `json:"ticket_id"`
	CustomerEmail string       `json:"customer_email"`
	Price         entity.Money `json:"price"`
}
