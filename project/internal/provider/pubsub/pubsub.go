package pubsub

import (
	"context"

	"github.com/ThreeDotsLabs/watermill/message"
)

type PubSub interface {
	message.Publisher

	AddConsumerHandler(
		handlerName string,
		subscribeTopic string,
		handlerFunc message.NoPublishHandlerFunc,
	) *message.Handler

	Register(ctx context.Context) error

	Running() chan struct{}
}
