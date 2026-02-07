package middleware

import (
	"errors"
	"tickets/internal/domain/errs"

	"github.com/ThreeDotsLabs/go-event-driven/v2/common/log"
	"github.com/ThreeDotsLabs/watermill/message"
)

func ErrorHandler(next message.HandlerFunc) message.HandlerFunc {
	return func(msg *message.Message) ([]*message.Message, error) {
		logger := log.FromContext(msg.Context())

		res, err := next(msg)
		if err == nil {
			return res, err
		}

		var appErr *errs.Error
		if errors.As(err, &appErr) && appErr.Code() != errs.CodeInternal {
			logger.With(
				"error", err,
				"error_code", appErr.Code(),
			).Error("Error while handling a message, skipping...")
			return nil, nil
		}

		logger.With("error", err).Error("Internal error while handling a message")
		return res, err
	}
}
