package entity

type Booking struct {
	ID              string `json:"id"`
	ShowID          string `json:"show_id"`
	NumberOfTickets int    `json:"number_of_tickets"`
	CustomerEmail   string `json:"customer_email"`
}
