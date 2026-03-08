package redisstream

import (
	"fmt"
	"os"
	"tickets/internal/pkg/bus"
	"tickets/internal/provider/cmdbus"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-redisstream/pkg/redisstream"
	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/redis/go-redis/v9"
)

const consumerGroupPrefix = "svc-tickets"

var marshaler = cqrs.JSONMarshaler{
	GenerateName: cqrs.StructName,
}

type CommandBus struct {
	processor *cqrs.CommandProcessor
	bus       *cqrs.CommandBus
	router    *message.Router
	pub       message.Publisher
}

func NewCommandBus() *CommandBus {
	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "localhost:6379"
	}

	rdb := redis.NewClient(&redis.Options{
		Addr: redisAddr,
	})

	var (
		pub message.Publisher
		err error
	)

	pub, err = redisstream.NewPublisher(redisstream.PublisherConfig{
		Client: rdb,
	}, nil)
	if err != nil {
		panic(err)
	}

	pub = bus.CorrelationPublisherDecorator{Publisher: pub}

	bus, err := cqrs.NewCommandBusWithConfig(
		pub,
		cqrs.CommandBusConfig{
			GeneratePublishTopic: generatePublishTopic,
			Marshaler:            marshaler,
		},
	)
	if err != nil {
		panic(err)
	}

	router := message.NewDefaultRouter(nil)

	processor, err := cqrs.NewCommandProcessorWithConfig(
		router,
		cqrs.CommandProcessorConfig{
			SubscriberConstructor:  newSubscriberConstructor(rdb, nil),
			Marshaler:              marshaler,
			GenerateSubscribeTopic: generateSubscribeTopic,
		},
	)
	if err != nil {
		panic(err)
	}

	return &CommandBus{
		bus:       bus,
		processor: processor,
		router:    router,
		pub:       pub,
	}
}

func newSubscriberConstructor(
	rdb *redis.Client,
	logger watermill.LoggerAdapter,
) cqrs.CommandProcessorSubscriberConstructorFn {
	return func(params cqrs.CommandProcessorSubscriberConstructorParams) (message.Subscriber, error) {
		return redisstream.NewSubscriber(redisstream.SubscriberConfig{
			Client:        rdb,
			ConsumerGroup: fmt.Sprintf("%s.%s", consumerGroupPrefix, params.HandlerName),
		}, logger)
	}
}

func generateSubscribeTopic(params cqrs.CommandProcessorGenerateSubscribeTopicParams) (string, error) {
	return params.CommandName, nil
}

func generatePublishTopic(params cqrs.CommandBusGeneratePublishTopicParams) (string, error) {
	return params.CommandName, nil
}

var _ cmdbus.CommandBus = (*CommandBus)(nil)
