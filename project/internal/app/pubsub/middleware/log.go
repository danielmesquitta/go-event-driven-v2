package middleware

import (
	"tickets/internal/pkg/ctxval"

	"tickets/internal/pkg/log"

	"github.com/ThreeDotsLabs/watermill/message"
)

func Logger(next message.HandlerFunc) message.HandlerFunc {
	return func(msg *message.Message) ([]*message.Message, error) {
		ctx := ctxval.WithMessageID(msg.Context(), msg.UUID)
		ctx = ctxval.WithPayload(ctx, msg.Payload)
		ctx = ctxval.WithMetadata(ctx, msg.Metadata)
		ctx = ctxval.WithHandlerName(ctx, message.HandlerNameFromCtx(ctx))
		msg.SetContext(ctx)

		logger := log.New(ctx)
		logger.Info("Handling message")

		res, err := next(msg)
		if err != nil {
			logger.With("error", err).Error("Error while handling a message")
		}

		return res, err
	}
}
