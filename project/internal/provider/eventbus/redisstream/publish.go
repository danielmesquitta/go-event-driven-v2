package redisstream

import (
	"context"
	"tickets/internal/app/pubsub/event"
)

func (e *EventBus) Publish(ctx context.Context, event event.Event) error {
	return e.bus.Publish(ctx, event)
}
