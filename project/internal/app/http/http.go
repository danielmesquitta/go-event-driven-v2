package http

import (
	"github.com/labstack/echo/v4"
)

func getCorrelationID(c echo.Context) string {
	return c.Request().Header.Get("Correlation-ID")
}
