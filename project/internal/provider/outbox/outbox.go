package outbox

import (
	"context"
	"tickets/internal/app/pubsub/message/event"
)

type Outbox interface {
	Publish(ctx context.Context, event event.Event) error
}
