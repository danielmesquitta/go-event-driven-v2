package http

import (
	"encoding/json"
	"net/http"
	"tickets/internal/app/pubsub/event"
	"tickets/internal/domain/entity"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/labstack/echo/v4"
)

type ticketsStatusRequest struct {
	Tickets []ticketStatusRequest `json:"tickets"`
}

type ticketStatusRequest struct {
	TicketID      string       `json:"ticket_id"`
	Status        string       `json:"status"`
	Price         entity.Money `json:"price"`
	CustomerEmail string       `json:"customer_email"`
}

func (h Handler) PostTicketsStatus(c echo.Context) error {
	var request ticketsStatusRequest
	err := c.Bind(&request)
	if err != nil {
		return err
	}

	for _, ticket := range request.Tickets {
		e := event.TicketBookingConfirmed{
			Header:        event.NewEventHeader(),
			TicketID:      ticket.TicketID,
			CustomerEmail: ticket.CustomerEmail,
			Price:         ticket.Price,
		}
		payload, err := json.Marshal(e)
		if err != nil {
			return err
		}
		msg := message.NewMessage(e.Header.ID, payload)
		if err := h.pubsub.Publish(event.TopicTicketBookingConfirmed, msg); err != nil {
			return err
		}
	}

	return c.NoContent(http.StatusOK)
}
