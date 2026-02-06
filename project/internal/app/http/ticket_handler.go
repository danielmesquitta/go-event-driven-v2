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

		if err := h.publishEvent(topic, e); err != nil {
			return err
		}
	}

	return c.NoContent(http.StatusOK)
}

func (h Handler) publishEvent(topic event.Topic, event event.Event) error {
	payload, err := json.Marshal(event)
	if err != nil {
		return err
	}
	msg := message.NewMessage(event.GetHeader().ID, payload)
	if err := h.pubsub.Publish(topic, msg); err != nil {
		return err
	}
	return nil
}
