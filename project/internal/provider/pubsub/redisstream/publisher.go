package redisstream

import (
	"tickets/internal/app/pubsub/event"

	"github.com/ThreeDotsLabs/watermill/message"
)

func (p *PubSub) Publish(topic event.Topic, msg *message.Message) error {
	return p.publisher.Publish(string(topic), msg)
}
