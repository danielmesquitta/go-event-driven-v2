package redisstream

import (
	"context"
	"fmt"

	"github.com/ThreeDotsLabs/watermill/message"
)

func (p *PubSub) AddConsumerHandler(
	handlerName string,
	subscribeTopic string,
	handlerFunc message.NoPublishHandlerFunc,
) *message.Handler {
	sub, err := p.newSubscriber(fmt.Sprintf("%s-%s-group", handlerName, subscribeTopic))
	if err != nil {
		panic(err)
	}

	return p.router.AddConsumerHandler(handlerName, subscribeTopic, sub, handlerFunc)
}

func (p *PubSub) Register(
	ctx context.Context,
) error {
	return p.router.Run(ctx)
}

func (p *PubSub) Running() chan struct{} {
	return p.router.Running()
}
