package pubsub

import (
	"context"
	"tickets/internal/app/pubsub/event"

	"github.com/ThreeDotsLabs/watermill/message"
)

type PubSub interface {
	Publish(topic event.Topic, msg *message.Message) error

	AddConsumerHandler(
		subscribeTopic event.Topic,
		handlerFunc message.NoPublishHandlerFunc,
	) *message.Handler

	AddHandler(
		subscribeTopic event.Topic,
		publishTopic event.Topic,
		handlerFunc message.HandlerFunc,
	) *message.Handler

	Register(ctx context.Context) error

	Running() chan struct{}
}
