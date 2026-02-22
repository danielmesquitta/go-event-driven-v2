package ctxval

import (
	"context"

	"github.com/google/uuid"
)

const IdempotencyKey ctxKey = "idempotency_key"

func WithIdempotencyKey(ctx context.Context, idempotencyKey string) context.Context {
	return context.WithValue(ctx, IdempotencyKey, idempotencyKey)
}

func GetIdempotencyKey(ctx context.Context) string {
	v, ok := ctx.Value(IdempotencyKey).(string)
	if ok {
		return v
	}
	return uuid.NewString()
}
