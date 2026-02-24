package handler

import (
	"net/http"
	"tickets/internal/domain/entity"
	"tickets/internal/domain/usecase"

	"github.com/labstack/echo/v4"
)

type TicketHandler struct {
	postTicketStatusUseCase *usecase.PostTicketStatus
	listTicketsUseCase      *usecase.ListTickets
}

func NewTicketHandler(
	postTicketStatusUseCase *usecase.PostTicketStatus,
	listTicketsUseCase *usecase.ListTickets,
) *TicketHandler {
	return &TicketHandler{
		postTicketStatusUseCase: postTicketStatusUseCase,
		listTicketsUseCase:      listTicketsUseCase,
	}
}

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

func (h TicketHandler) PostTicketsStatus(c echo.Context) error {
	var req ticketsStatusRequest
	err := c.Bind(&req)
	if err != nil {
		return err
	}

	tickets := make([]entity.Ticket, len(req.Tickets))
	for i, ticket := range req.Tickets {
		tickets[i] = entity.Ticket{
			ID:            ticket.TicketID,
			Status:        entity.TicketStatus(ticket.Status),
			Price:         ticket.Price,
			CustomerEmail: ticket.CustomerEmail,
		}
	}

	err = h.postTicketStatusUseCase.Execute(c.Request().Context(), usecase.PostTicketStatusInput{
		Tickets: tickets,
	})
	if err != nil {
		return err
	}

	return c.NoContent(http.StatusOK)
}

func (h TicketHandler) ListTickets(c echo.Context) error {
	tickets, err := h.listTicketsUseCase.Execute(c.Request().Context())
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, tickets)
}
