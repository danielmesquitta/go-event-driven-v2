package cmdbus

import (
	"context"
	"tickets/internal/app/pubsub/message/cmd"

	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/google/uuid"
)

type CommandBus interface {
	Send(ctx context.Context, cmd cmd.Command) error
	AddHandler(handler CommandHandler) error
	AddMiddleware(m ...message.HandlerMiddleware)

	Run(ctx context.Context) error

	// Running returns a channel that is closed when the pubsub system is started.
	Running() chan struct{}

	Router() *message.Router
	Publisher() message.Publisher
}

type HandlerFunc[T any] func(ctx context.Context, cmd *T) error

type CommandHandler interface {
	HandlerName() string
	NewCommand() any
	Handle(ctx context.Context, cmd any) error
}

func AddHandler[T any](eb CommandBus, handlerFunc HandlerFunc[T]) {
	handlerName := uuid.NewString()
	handler := cqrs.NewCommandHandler(handlerName, handlerFunc)
	err := eb.AddHandler(handler)
	if err != nil {
		panic(err)
	}
}
