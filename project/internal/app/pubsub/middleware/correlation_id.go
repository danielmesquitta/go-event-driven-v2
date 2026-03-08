package middleware

import (
	"tickets/internal/app/pubsub/message/event"
	"tickets/internal/pkg/ctxval"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/google/uuid"
)

func CorrelationID(next message.HandlerFunc) message.HandlerFunc {
	return func(msg *message.Message) ([]*message.Message, error) {
		correlationID := msg.Metadata.Get(string(event.MetadataKeyCorrelationID))
		if correlationID == "" {
			correlationID = uuid.NewString()
		}

		ctx := ctxval.WithCorrelationID(msg.Context(), correlationID)

		msg.SetContext(ctx)

		return next(msg)
	}
}
