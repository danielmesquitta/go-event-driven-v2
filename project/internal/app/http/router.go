package http

import (
	"net/http"
	"tickets/internal/provider/pubsub"

	libHttp "github.com/ThreeDotsLabs/go-event-driven/v2/common/http"
	"github.com/labstack/echo/v4"
)

func NewHttpRouter(
	pubsub pubsub.PubSub,
) *echo.Echo {
	e := libHttp.NewEcho()

	handler := Handler{
		pubsub: pubsub,
	}

	e.GET("/health", func(c echo.Context) error {
		return c.String(http.StatusOK, "ok")
	})

	e.POST("/tickets-status", handler.PostTicketsStatus)

	return e
}
