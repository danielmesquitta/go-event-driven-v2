package http

import (
	"net/http"
	"tickets/internal/app/pubsub/event"
	"tickets/internal/domain/entity"

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
		var e event.Event
		var topic event.Topic
		switch ticket.Status {
		case "confirmed":
			topic = event.TopicTicketBookingConfirmed
			e = &event.TicketBookingConfirmed{
				Header:        event.NewEventHeader(),
				TicketID:      ticket.TicketID,
				CustomerEmail: ticket.CustomerEmail,
				Price:         ticket.Price,
			}

		case "canceled":
			topic = event.TopicTicketBookingCanceled
			e = &event.TicketBookingCanceled{
				Header:        event.NewEventHeader(),
				TicketID:      ticket.TicketID,
				CustomerEmail: ticket.CustomerEmail,
				Price:         ticket.Price,
			}

		default:
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": "invalid status",
			})
		}

		if err := h.publishEvent(c, topic, e); err != nil {
			return err
		}
	}

	return c.NoContent(http.StatusOK)
}
