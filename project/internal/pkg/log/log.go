package log

import (
	"context"
	"log/slog"
	"tickets/internal/pkg/ctxval"
)

func New(ctx context.Context) *slog.Logger {
	logger := slog.Default()

	correlationID := ctxval.GetCorrelationID(ctx)
	if correlationID != "" {
		logger = logger.With(string(ctxval.CorrelationIDKey), correlationID)
	}

	idempotencyKey := ctxval.GetIdempotencyKey(ctx)
	if idempotencyKey != "" {
		logger = logger.With(string(ctxval.IdempotencyKey), idempotencyKey)
	}

	messageID := ctxval.GetMessageID(ctx)
	if messageID != "" {
		logger = logger.With(string(ctxval.MessageIDKey), messageID)
	}

	payload := ctxval.GetPayload(ctx)
	if len(payload) > 0 {
		logger = logger.With(string(ctxval.PayloadKey), string(payload))
	}

	metadata := ctxval.GetMetadata(ctx)
	if len(metadata) > 0 {
		logger = logger.With(string(ctxval.MetadataKey), metadata)
	}

	handlerName := ctxval.GetHandlerName(ctx)
	if handlerName != "" {
		logger = logger.With(string(ctxval.HandlerNameKey), handlerName)
	}

	return logger
}
