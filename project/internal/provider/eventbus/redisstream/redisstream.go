package redisstream

import (
	"fmt"
	"os"
	"tickets/internal/provider/eventbus"

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

type EventBus struct {
	processor *cqrs.EventProcessor
	bus       *cqrs.EventBus
	router    *message.Router
}

func NewEventBus() *EventBus {
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

	pub = eventbus.CorrelationPublisherDecorator{Publisher: pub}

	bus, err := cqrs.NewEventBusWithConfig(
		pub,
		cqrs.EventBusConfig{
			GeneratePublishTopic: generatePublishTopic,
			Marshaler:            marshaler,
		},
	)
	if err != nil {
		panic(err)
	}

	router := message.NewDefaultRouter(nil)

	processor, err := cqrs.NewEventProcessorWithConfig(
		router,
		cqrs.EventProcessorConfig{
			SubscriberConstructor:  newSubscriberConstructor(rdb, nil),
			Marshaler:              marshaler,
			GenerateSubscribeTopic: generateSubscribeTopic,
		},
	)
	if err != nil {
		panic(err)
	}

	return &EventBus{
		bus:       bus,
		processor: processor,
		router:    router,
	}
}

func newSubscriberConstructor(
	rdb *redis.Client,
	logger watermill.LoggerAdapter,
) cqrs.EventProcessorSubscriberConstructorFn {
	return func(params cqrs.EventProcessorSubscriberConstructorParams) (message.Subscriber, error) {
		return redisstream.NewSubscriber(redisstream.SubscriberConfig{
			Client:        rdb,
			ConsumerGroup: fmt.Sprintf("%s.%s", consumerGroupPrefix, params.HandlerName),
		}, logger)
	}
}

func generateSubscribeTopic(params cqrs.EventProcessorGenerateSubscribeTopicParams) (string, error) {
	return params.EventName, nil
}

func generatePublishTopic(params cqrs.GenerateEventPublishTopicParams) (string, error) {
	return params.EventName, nil
}

var _ eventbus.EventBus = (*EventBus)(nil)
