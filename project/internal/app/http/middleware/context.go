package middleware

import (
	"net/http"
	"tickets/internal/pkg/ctxval"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

const (
	HeaderIdempotencyKey = "Idempotency-Key"
	HeaderCorrelationID  = "Correlation-ID"
)

func Context(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {

		handleCorrelationID(c)

		err := handleIdempotencyKey(c)
		if err != nil {
			return err
		}

		return next(c)
	}
}

func handleCorrelationID(c echo.Context) {
	ctx := c.Request().Context()

	correlationID := c.Request().Header.Get(HeaderCorrelationID)
	if correlationID == "" {
		correlationID = uuid.NewString()
		c.Response().Header().Set(HeaderCorrelationID, correlationID)
	}

	ctx = ctxval.WithCorrelationID(ctx, correlationID)
	c.SetRequest(c.Request().WithContext(ctx))
}

func handleIdempotencyKey(c echo.Context) error {
	method := c.Request().Method
	if method != http.MethodPost {
		return nil
	}

	ctx := c.Request().Context()

	idempotencyKey := c.Request().Header.Get(HeaderIdempotencyKey)
	if idempotencyKey == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "idempotency key is required")
	}

	ctx = ctxval.WithIdempotencyKey(ctx, idempotencyKey)
	c.SetRequest(c.Request().WithContext(ctx))

	return nil
}
