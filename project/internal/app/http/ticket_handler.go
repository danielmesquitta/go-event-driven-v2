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

type ticketStatus string

const (
	statusConfirmed ticketStatus = "confirmed"
	statusCanceled  ticketStatus = "canceled"
)

type ticketStatusRequest struct {
	TicketID      string       `json:"ticket_id"`
	Status        ticketStatus `json:"status"`
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
		switch ticket.Status {
		case statusConfirmed:
			e = &event.TicketBookingConfirmed{
				Header:        event.NewEventHeader(),
				TicketID:      ticket.TicketID,
				CustomerEmail: ticket.CustomerEmail,
				Price:         ticket.Price,
			}

		case statusCanceled:
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

		if err := publishEvent(c, h.eventBus, e); err != nil {
			return err
		}
	}

	return c.NoContent(http.StatusOK)
}
