package redisstream

import "github.com/ThreeDotsLabs/watermill/message"

func (e *CommandBus) AddMiddleware(m ...message.HandlerMiddleware) {
	e.router.AddMiddleware(m...)
}
