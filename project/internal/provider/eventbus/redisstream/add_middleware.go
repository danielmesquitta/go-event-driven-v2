package redisstream

import "github.com/ThreeDotsLabs/watermill/message"

func (e *EventBus) AddMiddleware(m ...message.HandlerMiddleware) {
	e.router.AddMiddleware(m...)
}
