package http

import (
	"net/http"
	"tickets/internal/provider/eventbus"

	libHttp "github.com/ThreeDotsLabs/go-event-driven/v2/common/http"
	"github.com/labstack/echo/v4"
)

func NewHttpRouter(
	eventBus eventbus.EventBus,
) *echo.Echo {
	e := libHttp.NewEcho()

	handler := Handler{
		eventBus: eventBus,
	}

	e.GET("/health", func(c echo.Context) error {
		return c.String(http.StatusOK, "ok")
	})

	e.POST("/tickets-status", handler.PostTicketsStatus)

	return e
}
