package http

import (
	"net/http"
	"tickets/internal/app/pubsub/event"
	"tickets/internal/domain/entity"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/labstack/echo/v4"
)

type ticketsConfirmationRequest struct {
	Tickets []string `json:"tickets"`
}

func (h Handler) PostTicketsConfirmation(c echo.Context) error {
	var request ticketsConfirmationRequest
	err := c.Bind(&request)
	if err != nil {
		return err
	}

	for _, ticket := range request.Tickets {
		msg := message.NewMessage(watermill.NewUUID(), []byte(ticket))
		if err := h.pubsub.Publish(event.TopicIssueReceipt, msg); err != nil {
			return err
		}

		msg = message.NewMessage(watermill.NewUUID(), []byte(ticket))
		if err := h.pubsub.Publish(event.TopicAppendToTracker, msg); err != nil {
			return err
		}
	}

	return c.NoContent(http.StatusOK)
}

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
		msg := message.NewMessage(watermill.NewUUID(), []byte(ticket.TicketID))
		if err := h.pubsub.Publish(event.TopicIssueReceipt, msg); err != nil {
			return err
		}

		msg = message.NewMessage(watermill.NewUUID(), []byte(ticket.TicketID))
		if err := h.pubsub.Publish(event.TopicAppendToTracker, msg); err != nil {
			return err
		}
	}

	return c.NoContent(http.StatusOK)
}
