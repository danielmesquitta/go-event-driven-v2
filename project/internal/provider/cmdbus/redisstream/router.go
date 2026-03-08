package redisstream

import (
	"context"

	"github.com/ThreeDotsLabs/watermill/message"
)

func (e *CommandBus) Run(ctx context.Context) error {
	return e.router.Run(ctx)
}

func (e *CommandBus) Running() chan struct{} {
	return e.router.Running()
}

func (e *CommandBus) Router() *message.Router {
	return e.router
}

func (e *CommandBus) Publisher() message.Publisher {
	return e.pub
}
