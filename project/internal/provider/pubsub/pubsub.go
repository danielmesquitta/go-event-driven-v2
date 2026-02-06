package pubsub

import (
	"context"

	"github.com/ThreeDotsLabs/watermill/message"
)

type PubSub interface {
	message.Publisher

	NewSubscriber(consumerGroup string) (message.Subscriber, error)

	Register(ctx context.Context) error

	Running() chan struct{}
}
