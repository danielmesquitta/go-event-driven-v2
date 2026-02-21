package redisstream

import (
	"context"
	"tickets/internal/app/pubsub/event"
	"tickets/internal/provider/eventbus"
)

func (e *EventBus) Publish(ctx context.Context, ev event.Event, opts ...eventbus.PublishOption) error {
	ctx = eventbus.ApplyPublishOptions(ctx, opts)
	return e.bus.Publish(ctx, ev)
}
