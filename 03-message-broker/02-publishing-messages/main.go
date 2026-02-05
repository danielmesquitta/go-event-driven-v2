package main

import (
	"os"
	"strconv"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-redisstream/pkg/redisstream"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

func main() {
	logger := watermill.NewSlogLogger(nil)

	rdb := redis.NewClient(&redis.Options{
		Addr: os.Getenv("REDIS_ADDR"),
	})

	publisher, err := redisstream.NewPublisher(redisstream.PublisherConfig{
		Client: rdb,
	}, logger)
	if err != nil {
		panic(err)
	}

	for _, progress := range []int{50, 100} {
		msg := message.NewMessage(uuid.NewString(), []byte(strconv.Itoa(progress)))
		err = publisher.Publish("progress", msg)
		if err != nil {
			panic(err)
		}
	}
}
