package redisstream

import (
	"context"
	"fmt"
	"tickets/internal/provider/pubsub"

	"github.com/ThreeDotsLabs/watermill/message"
)

func (p *PubSub) AddConsumerHandler(
	handlerName string,
	subscribeTopic string,
	handlerFunc message.NoPublishHandlerFunc,
) *message.Handler {
	sub, err := p.NewSubscriber(fmt.Sprintf("%s-%s-group", handlerName, subscribeTopic))
	if err != nil {
		panic(err)
	}

	return p.router.AddConsumerHandler(handlerName, subscribeTopic, sub, handlerFunc)
}

func (p *PubSub) Register(
	ctx context.Context,
) error {
	p.AddConsumerHandler(
		"HandleIssueReceipt",
		pubsub.TopicIssueReceipt,
		p.handleIssueReceipt,
	)

	p.AddConsumerHandler(
		"HandleAppendToTracker",
		pubsub.TopicAppendToTracker,
		p.handleAppendToTracker,
	)

	return p.router.Run(ctx)
}

func (p *PubSub) Running() chan struct{} {
	return p.router.Running()
}
