package outbox

import (
	"context"
	"tickets/internal/app/pubsub/event"
)

type Outbox interface {
	Publish(ctx context.Context, event event.Event) error
}
