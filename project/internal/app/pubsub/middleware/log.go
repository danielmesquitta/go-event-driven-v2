package middleware

import (
	"github.com/ThreeDotsLabs/go-event-driven/v2/common/log"
	"github.com/ThreeDotsLabs/watermill/message"
)

func Logger(next message.HandlerFunc) message.HandlerFunc {
	return func(msg *message.Message) ([]*message.Message, error) {
		logger := log.FromContext(msg.Context())
		logger = logger.With(
			"message_id", msg.UUID,
			"payload", string(msg.Payload),
			"metadata", msg.Metadata,
			"handler", message.HandlerNameFromCtx(msg.Context()),
		)

		ctx := log.ToContext(msg.Context(), logger)
		msg.SetContext(ctx)

		logger.Info("Handling a message")

		res, err := next(msg)
		if err != nil {
			logger.With("error", err).Error("Error while handling a message")
		}

		return res, err
	}
}
