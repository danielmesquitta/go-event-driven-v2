package ctxval

import (
	"context"

	"github.com/google/uuid"
)

const (
	MessageIDKey   ctxKey = "message_id"
	PayloadKey     ctxKey = "payload"
	MetadataKey    ctxKey = "metadata"
	HandlerNameKey ctxKey = "handler_name"
)

func WithMessageID(ctx context.Context, messageID string) context.Context {
	return context.WithValue(ctx, MessageIDKey, messageID)
}

func GetMessageID(ctx context.Context) string {
	v, ok := ctx.Value(MessageIDKey).(string)
	if ok {
		return v
	}
	return uuid.NewString()
}

func WithPayload(ctx context.Context, payload []byte) context.Context {
	return context.WithValue(ctx, PayloadKey, payload)
}

func GetPayload(ctx context.Context) []byte {
	v, ok := ctx.Value(PayloadKey).([]byte)
	if ok {
		return v
	}
	return nil
}

func WithMetadata(ctx context.Context, metadata map[string]string) context.Context {
	return context.WithValue(ctx, MetadataKey, metadata)
}

func GetMetadata(ctx context.Context) map[string]string {
	v, ok := ctx.Value(MetadataKey).(map[string]string)
	if ok {
		return v
	}
	return make(map[string]string)
}

func WithHandlerName(ctx context.Context, handlerName string) context.Context {
	return context.WithValue(ctx, HandlerNameKey, handlerName)
}

func GetHandlerName(ctx context.Context) string {
	v, ok := ctx.Value(HandlerNameKey).(string)
	if ok {
		return v
	}
	return ""
}
