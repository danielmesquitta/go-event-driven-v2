package http

import (
	"github.com/ThreeDotsLabs/watermill/message"

	libHttp "github.com/ThreeDotsLabs/go-event-driven/v2/common/http"
	"github.com/labstack/echo/v4"
)

func NewHttpRouter(
	publisher message.Publisher,
) *echo.Echo {
	e := libHttp.NewEcho()

	handler := Handler{
		publisher: publisher,
	}

	e.POST("/tickets-confirmation", handler.PostTicketsConfirmation)

	return e
}
