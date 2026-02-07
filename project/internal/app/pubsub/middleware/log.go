package middleware

import (
	"log/slog"

	"github.com/ThreeDotsLabs/watermill/message"
)

func Logger(next message.HandlerFunc) message.HandlerFunc {
	return func(msg *message.Message) ([]*message.Message, error) {
		logger := slog.With(
			"message_id", msg.UUID,
			"payload", string(msg.Payload),
			"metadata", msg.Metadata,
			"handler", message.HandlerNameFromCtx(msg.Context()),
		)

		logger.Info("Handling a message")

		return next(msg)
	}
}
