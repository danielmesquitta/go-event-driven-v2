package redisstream

import (
	"context"
)

func (p *PubSub) Run(
	ctx context.Context,
) error {
	return p.router.Run(ctx)
}

func (p *PubSub) Running() chan struct{} {
	return p.router.Running()
}
