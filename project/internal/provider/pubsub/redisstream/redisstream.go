package redisstream

import (
	"os"
	"tickets/internal/provider/pubsub"

	"github.com/ThreeDotsLabs/watermill-redisstream/pkg/redisstream"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/redis/go-redis/v9"
)

type PubSub struct {
	publisher message.Publisher

	router *message.Router
	rdb    *redis.Client
}

func NewPubSub() *PubSub {
	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "localhost:6379"
	}

	rdb := redis.NewClient(&redis.Options{
		Addr: redisAddr,
	})

	publisher, err := newRedisPublisher(rdb)
	if err != nil {
		panic(err)
	}

	return &PubSub{
		publisher: publisher,
		router:    message.NewDefaultRouter(nil),
		rdb:       rdb,
	}
}

func newRedisPublisher(rdb *redis.Client) (message.Publisher, error) {
	return redisstream.NewPublisher(redisstream.PublisherConfig{
		Client: rdb,
	}, nil)
}

func (p *PubSub) newSubscriber(consumerGroup string) (message.Subscriber, error) {
	return redisstream.NewSubscriber(redisstream.SubscriberConfig{
		Client:        p.rdb,
		ConsumerGroup: consumerGroup,
	}, nil)
}

var _ pubsub.PubSub = (*PubSub)(nil)
