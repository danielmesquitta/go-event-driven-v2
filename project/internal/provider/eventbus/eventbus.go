package eventbus

import (
	"context"
	"tickets/internal/app/pubsub/event"

	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/google/uuid"
)

type PublishOption func(*publishOptions)

type publishOptions struct {
	correlationID string
}

func WithCorrelationID(id string) PublishOption {
	return func(o *publishOptions) {
		o.correlationID = id
	}
}

func ApplyPublishOptions(ctx context.Context, opts []PublishOption) context.Context {
	o := &publishOptions{}
	for _, opt := range opts {
		opt(o)
	}
	if o.correlationID != "" {
		ctx = ContextWithCorrelationID(ctx, o.correlationID)
	}
	return ctx
}

type EventBus interface {
	Publish(ctx context.Context, event event.Event, opts ...PublishOption) error
	AddHandler(handler EventHandler) error
	AddMiddleware(m ...message.HandlerMiddleware)

	Run(ctx context.Context) error

	// Running returns a channel that is closed when the pubsub system is started.
	Running() chan struct{}
}

type HandlerFunc[T any] func(ctx context.Context, event *T) error

type EventHandler interface {
	HandlerName() string
	NewEvent() any
	Handle(ctx context.Context, event any) error
}

func AddHandler[T any](eb EventBus, handlerFunc HandlerFunc[T]) {
	handlerName := uuid.NewString()
	handler := cqrs.NewEventHandler(handlerName, handlerFunc)
	err := eb.AddHandler(handler)
	if err != nil {
		panic(err)
	}
}
