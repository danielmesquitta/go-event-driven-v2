package middleware

import (
	"net/http"
	"tickets/internal/pkg/ctxval"

	"github.com/google/uuid"
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
			// Should return error, but for development purposes we will create a new idempotency key
			// return echo.NewHTTPError(http.StatusBadRequest, "idempotency key is required")
			idempotencyKey = uuid.NewString()
			c.Response().Header().Set(HeaderIdempotencyKey, idempotencyKey)
		}

		ctx = ctxval.WithIdempotencyKey(ctx, idempotencyKey)
		c.SetRequest(c.Request().WithContext(ctx))

		return next(c)
	}
}
