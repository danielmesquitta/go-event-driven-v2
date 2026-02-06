package pubsub

import (
	"context"
	"tickets/internal/app/pubsub/event"

	"github.com/ThreeDotsLabs/watermill/message"
)

type PubSub interface {
	Publish(topic event.Topic, msg *message.Message) error

	// AddConsumerHandler listens to a topic and calls the handler function for each message.
	// It creates a new consumer group for each handler function.
	AddConsumerHandler(
		subscribeTopic event.Topic,
		handlerFunc message.NoPublishHandlerFunc,
	) *message.Handler

	// AddHandler listens to the subscribe topic and sends the message(s) to the publish topic.
	// It creates a new consumer group for each handler function.
	AddHandler(
		subscribeTopic event.Topic,
		publishTopic event.Topic,
		handlerFunc message.HandlerFunc,
	) *message.Handler

	// Run starts the pubsub system, with all the handlers and consumer handlers registered.
	// It blocks until the context is cancelled.
	Run(ctx context.Context) error

	// Running returns a channel that is closed when the pubsub system is started.
	Running() chan struct{}
}
