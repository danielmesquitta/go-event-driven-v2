package ctxval

import (
	"context"
)

const CorrelationIDKey ctxKey = "correlation_id"

func WithCorrelationID(ctx context.Context, correlationID string) context.Context {
	return context.WithValue(ctx, CorrelationIDKey, correlationID)
}

func GetCorrelationID(ctx context.Context) string {
	v, ok := ctx.Value(CorrelationIDKey).(string)
	if ok {
		return v
	}
	return ""
}
