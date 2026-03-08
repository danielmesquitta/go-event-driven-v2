package redisstream

import (
	"context"
	"tickets/internal/app/pubsub/message/cmd"
)

func (e *CommandBus) Send(ctx context.Context, ev cmd.Command) error {
	return e.bus.Send(ctx, ev)
}
