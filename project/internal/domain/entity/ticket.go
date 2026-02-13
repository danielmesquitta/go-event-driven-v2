package entity

type Ticket struct {
	ID            string `json:"ticket_id"`
	Price         Money  `json:"price"`
	CustomerEmail string `json:"customer_email"`
}
