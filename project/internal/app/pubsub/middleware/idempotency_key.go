package middleware

import (
	"encoding/json"
	pubSubMessage "tickets/internal/app/pubsub/message"
	"tickets/internal/pkg/ctxval"

	"github.com/ThreeDotsLabs/watermill/message"
)

func IdempotencyKey(next message.HandlerFunc) message.HandlerFunc {
	return func(msg *message.Message) ([]*message.Message, error) {
		type Event struct {
			Header pubSubMessage.Header `json:"header"`
		}

		e := Event{}
		err := json.Unmarshal(msg.Payload, &e)
		if err != nil {
			return nil, err
		}

		ctx := ctxval.WithIdempotencyKey(msg.Context(), e.Header.IdempotencyKey)
		msg.SetContext(ctx)

		return next(msg)
	}
}
