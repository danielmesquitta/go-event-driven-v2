package redisstream

import "github.com/ThreeDotsLabs/watermill/message"

func (p *PubSub) AddMiddleware(m ...message.HandlerMiddleware) {
	p.router.AddMiddleware(m...)
}
