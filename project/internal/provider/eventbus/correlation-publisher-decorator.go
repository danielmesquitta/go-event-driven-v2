package eventbus

import (
	"tickets/internal/pkg/ctxval"

	"github.com/ThreeDotsLabs/watermill/message"
)

type CorrelationPublisherDecorator struct {
	message.Publisher
}

func (c CorrelationPublisherDecorator) Publish(topic string, messages ...*message.Message) error {
	for _, message := range messages {
		message.Metadata.Set("correlation_id", ctxval.GetCorrelationID(message.Context()))
	}

	return c.Publisher.Publish(topic, messages...)
}
