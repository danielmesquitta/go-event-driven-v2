package event

import "tickets/internal/app/pubsub/message"

type Event interface {
	// GetHeader returns the header of the event.
	GetHeader() message.Header
}
