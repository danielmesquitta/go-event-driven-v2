package pubsub

import (
	"context"
	"tickets/internal/app/pubsub/event"

	"github.com/ThreeDotsLabs/watermill/message"
)

type PubSub interface {
	Publish(topic event.Topic, msg *message.Message) error

	AddConsumerHandler(
		handlerName string,
		subscribeTopic event.Topic,
		handlerFunc message.NoPublishHandlerFunc,
	) *message.Handler

	Register(ctx context.Context) error

	Running() chan struct{}
}
