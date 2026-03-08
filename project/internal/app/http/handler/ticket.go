package handler

import (
	"net/http"
	"tickets/internal/domain/entity"
	"tickets/internal/domain/usecase/ticket"

	"github.com/labstack/echo/v4"
)

type TicketHandler struct {
	postTicketStatusUseCase        *ticket.PostStatus
	listTicketsUseCase             *ticket.List
	sendTicketRefundCommandUseCase *ticket.SendRefundCommand
}

func NewTicketHandler(
	postTicketStatusUseCase *ticket.PostStatus,
	listTicketsUseCase *ticket.List,
	sendTicketRefundCommandUseCase *ticket.SendRefundCommand,
) *TicketHandler {
	return &TicketHandler{
		postTicketStatusUseCase:        postTicketStatusUseCase,
		listTicketsUseCase:             listTicketsUseCase,
		sendTicketRefundCommandUseCase: sendTicketRefundCommandUseCase,
	}
}

type ticketsStatusRequest struct {
	Tickets []ticketStatusRequest `json:"tickets"`
}

type ticketStatusRequest struct {
	TicketID      string       `json:"ticket_id"`
	BookingID     string       `json:"booking_id"`
	Status        string       `json:"status"`
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
	for i, t := range req.Tickets {
		tickets[i] = entity.Ticket{
			ID:            t.TicketID,
			Status:        entity.TicketStatus(t.Status),
			Price:         t.Price,
			CustomerEmail: t.CustomerEmail,
			BookingID:     t.BookingID,
		}
	}

	err = h.postTicketStatusUseCase.Execute(c.Request().Context(), ticket.PostStatusInput{
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

func (h TicketHandler) TicketRefund(c echo.Context) error {
	ticketID := c.Param("ticket_id")

	err := h.sendTicketRefundCommandUseCase.Execute(c.Request().Context(), ticket.SendRefundCommandInput{
		TicketID: ticketID,
	})
	if err != nil {
		return err
	}

	return c.NoContent(http.StatusAccepted)
}
