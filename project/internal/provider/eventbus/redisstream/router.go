package redisstream

import (
	"context"

	"github.com/ThreeDotsLabs/watermill/message"
)

func (e *EventBus) Run(ctx context.Context) error {
	return e.router.Run(ctx)
}

func (e *EventBus) Running() chan struct{} {
	return e.router.Running()
}

func (e *EventBus) Router() *message.Router {
	return e.router
}

func (e *EventBus) Publisher() message.Publisher {
	return e.pub
}
