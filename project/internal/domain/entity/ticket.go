package entity

type TicketStatus string

const (
	TicketStatusConfirmed TicketStatus = "confirmed"
	TicketStatusCanceled  TicketStatus = "canceled"
)

type Ticket struct {
	ID            string       `json:"ticket_id"`
	Status        TicketStatus `json:"status" validate:"oneof=confirmed canceled"`
	Price         Money        `json:"price"`
	CustomerEmail string       `json:"customer_email"`
}
