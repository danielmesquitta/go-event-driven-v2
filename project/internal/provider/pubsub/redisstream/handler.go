package redisstream

import (
	"fmt"
	"tickets/internal/app/pubsub/event"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/google/uuid"
)

func (p *PubSub) AddHandler(
	subscribeTopic event.Topic,
	publishTopic event.Topic,
	handlerFuncs ...message.HandlerFunc,
) {
	for _, handlerFunc := range handlerFuncs {
		id := uuid.NewString()

		sub, err := p.newSubscriber(fmt.Sprintf("%s-%s-group", id, subscribeTopic))
		if err != nil {
			panic(err)
		}

		p.router.AddHandler(id, string(subscribeTopic), sub, string(publishTopic), p.publisher, handlerFunc)
	}
}
