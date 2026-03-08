package middleware

import (
	"errors"
	"net/http"
	"tickets/internal/domain/errs"
	"tickets/internal/pkg/log"

	"github.com/labstack/echo/v4"
)

type errorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func ErrorHandler(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		err := next(c)
		if err == nil {
			return nil
		}

		logger := log.New(c.Request().Context())

		var appErr *errs.Error
		if errors.As(err, &appErr) {
			status := httpStatusFromCode(appErr.Code())

			logger.With(
				"error", err,
				"error_code", appErr.Code(),
				"http_status", status,
			).Error("Domain error in HTTP handler")

			return c.JSON(status, errorResponse{
				Code:    string(appErr.Code()),
				Message: appErr.Error(),
			})
		}

		var echoErr *echo.HTTPError
		if errors.As(err, &echoErr) {
			return err
		}

		logger.With("error", err).Error("Internal error in HTTP handler")
		return c.JSON(http.StatusInternalServerError, errorResponse{
			Code:    string(errs.CodeInternal),
			Message: "internal server error",
		})
	}
}

func httpStatusFromCode(code errs.Code) int {
	switch code {
	case errs.CodeBadRequest:
		return http.StatusBadRequest
	case errs.CodeNotFound:
		return http.StatusNotFound
	case errs.CodeUnauthorized:
		return http.StatusUnauthorized
	case errs.CodeForbidden:
		return http.StatusForbidden
	default:
		return http.StatusInternalServerError
	}
}
