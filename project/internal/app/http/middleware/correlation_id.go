package middleware

import (
	"tickets/internal/pkg/ctxval"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

const (
	HeaderCorrelationID = "Correlation-ID"
)

func CorrelationID(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()

		correlationID := c.Request().Header.Get(HeaderCorrelationID)
		if correlationID == "" {
			correlationID = uuid.NewString()
			c.Response().Header().Set(HeaderCorrelationID, correlationID)
		}

		ctx = ctxval.WithCorrelationID(ctx, correlationID)
		c.SetRequest(c.Request().WithContext(ctx))

		return next(c)
	}
}
