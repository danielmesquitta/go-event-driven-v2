package http

import (
	"net/http"
	"tickets/internal/provider/eventbus"
	"tickets/internal/provider/repo"

	libHttp "github.com/ThreeDotsLabs/go-event-driven/v2/common/http"
	"github.com/labstack/echo/v4"
)

func NewHttpRouter(
	eventBus eventbus.EventBus,
	ticketRepo repo.TicketRepo,
) *echo.Echo {
	e := libHttp.NewEcho()

	handler := NewHandler(eventBus, ticketRepo)

	e.GET("/health", func(c echo.Context) error {
		return c.String(http.StatusOK, "ok")
	})

	e.POST("/tickets-status", handler.PostTicketsStatus)
	e.GET("/tickets", handler.ListTickets)

	return e
}
