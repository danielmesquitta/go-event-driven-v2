package redisstream

import "context"

func (e *EventBus) Run(ctx context.Context) error {
	return e.router.Run(ctx)
}

func (e *EventBus) Running() chan struct{} {
	return e.router.Running()
}
