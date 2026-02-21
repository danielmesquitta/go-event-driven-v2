package http

import (
	"net/http"
	"tickets/internal/domain/usecase"

	libHttp "github.com/ThreeDotsLabs/go-event-driven/v2/common/http"
	"github.com/labstack/echo/v4"
)

func NewHttpRouter(
	postTicketStatusUseCase *usecase.PostTicketStatus,
	listTicketsUseCase *usecase.ListTickets,
) *echo.Echo {
	e := libHttp.NewEcho()

	ticketHandler := NewTicketHandler(
		postTicketStatusUseCase,
		listTicketsUseCase,
	)

	e.GET("/health", func(c echo.Context) error {
		return c.String(http.StatusOK, "ok")
	})

	e.POST("/tickets-status", ticketHandler.PostTicketsStatus)
	e.GET("/tickets", ticketHandler.ListTickets)

	return e
}
