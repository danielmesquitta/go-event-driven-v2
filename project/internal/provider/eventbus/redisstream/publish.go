package redisstream

import (
	"context"
	"tickets/internal/app/pubsub/message/event"
)

func (e *EventBus) Publish(ctx context.Context, ev event.Event) error {
	return e.bus.Publish(ctx, ev)
}
