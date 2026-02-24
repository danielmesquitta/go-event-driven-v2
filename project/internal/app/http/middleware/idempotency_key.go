package middleware

import (
	"net/http"
	"tickets/internal/pkg/ctxval"

	"github.com/labstack/echo/v4"
)

const (
	HeaderIdempotencyKey = "Idempotency-Key"
)

func IdempotencyKey(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		method := c.Request().Method
		if method != http.MethodPost {
			return next(c)
		}

		ctx := c.Request().Context()

		idempotencyKey := c.Request().Header.Get(HeaderIdempotencyKey)
		if idempotencyKey == "" {
			return echo.NewHTTPError(http.StatusBadRequest, "idempotency key is required")
		}

		ctx = ctxval.WithIdempotencyKey(ctx, idempotencyKey)
		c.SetRequest(c.Request().WithContext(ctx))

		return next(c)
	}
}
