package middleware

import (
	"log/slog"

	"github.com/ThreeDotsLabs/go-event-driven/v2/common/log"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/google/uuid"
)

func CorrelationID(next message.HandlerFunc) message.HandlerFunc {
	return func(msg *message.Message) ([]*message.Message, error) {
		correlationID := msg.Metadata.Get("correlation_id")
		if correlationID == "" {
			correlationID = uuid.NewString()
		}

		ctx := log.ContextWithCorrelationID(msg.Context(), correlationID)
		ctx = log.ToContext(ctx, slog.With("correlation_id", correlationID))

		msg.SetContext(ctx)

		return next(msg)
	}
}
