package redisstream

import (
	"fmt"
	"tickets/internal/app/pubsub/event"

	"github.com/ThreeDotsLabs/watermill/message"
)

func (p *PubSub) AddConsumerHandler(
	subscribeTopic event.Topic,
	handlerFunc message.NoPublishHandlerFunc,
) *message.Handler {
	handlerName := handlerNameFromFunc(handlerFunc)

	sub, err := p.newSubscriber(fmt.Sprintf("%s-%s-group", handlerName, subscribeTopic))
	if err != nil {
		panic(err)
	}

	return p.router.AddConsumerHandler(handlerName, string(subscribeTopic), sub, handlerFunc)
}
