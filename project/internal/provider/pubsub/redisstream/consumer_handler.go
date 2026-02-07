package redisstream

import (
	"fmt"
	"tickets/internal/app/pubsub/event"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/google/uuid"
)

func (p *PubSub) AddConsumerHandler(
	subscribeTopic event.Topic,
	handlerFuncs ...message.NoPublishHandlerFunc,
) {
	for _, handlerFunc := range handlerFuncs {
		id := uuid.NewString()

		sub, err := p.newSubscriber(fmt.Sprintf("%s-%s-group", id, subscribeTopic))
		if err != nil {
			panic(err)
		}

		p.router.AddConsumerHandler(id, string(subscribeTopic), sub, handlerFunc)
	}
}
