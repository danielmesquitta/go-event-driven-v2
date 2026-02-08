package http

import (
	"tickets/internal/app/pubsub/event"
	"tickets/internal/provider/eventbus"

	"github.com/labstack/echo/v4"
)

type Handler struct {
	eventBus eventbus.EventBus
}

func getCorrelationID(c echo.Context) string {
	return c.Request().Header.Get("Correlation-ID")
}

func publishEvent(c echo.Context, eventBus eventbus.EventBus, event event.Event) error {
	ctx := c.Request().Context()

	correlationID := getCorrelationID(c)
	if correlationID != "" {
		ctx = eventbus.ContextWithCorrelationID(ctx, correlationID)
	}

	return eventBus.Publish(ctx, event)
}
