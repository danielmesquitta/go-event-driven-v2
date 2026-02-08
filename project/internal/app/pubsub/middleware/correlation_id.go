package middleware

import (
	"log/slog"
	"tickets/internal/app/pubsub/event"
	"tickets/internal/provider/eventbus"

	"github.com/ThreeDotsLabs/go-event-driven/v2/common/log"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/google/uuid"
)

func CorrelationID(next message.HandlerFunc) message.HandlerFunc {
	return func(msg *message.Message) ([]*message.Message, error) {
		correlationID := msg.Metadata.Get(string(event.MetadataKeyCorrelationID))
		if correlationID == "" {
			correlationID = uuid.NewString()
		}

		ctx := eventbus.ContextWithCorrelationID(msg.Context(), correlationID)
		ctx = log.ToContext(ctx, slog.With(string(event.MetadataKeyCorrelationID), correlationID))

		msg.SetContext(ctx)

		return next(msg)
	}
}
