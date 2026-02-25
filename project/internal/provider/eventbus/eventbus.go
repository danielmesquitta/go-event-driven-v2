package eventbus

import (
	"context"
	"tickets/internal/app/pubsub/event"

	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/google/uuid"
)

type EventBus interface {
	Publish(ctx context.Context, event event.Event) error
	AddHandler(handler EventHandler) error
	AddMiddleware(m ...message.HandlerMiddleware)

	Run(ctx context.Context) error

	// Running returns a channel that is closed when the pubsub system is started.
	Running() chan struct{}

	Router() *message.Router
	Publisher() message.Publisher
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
