package entity

type Ticket struct {
	ID            string `json:"id"`
	Price         Money  `json:"price"`
	CustomerEmail string `json:"customer_email"`
}
