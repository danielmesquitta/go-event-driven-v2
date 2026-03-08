package cmd

import "tickets/internal/app/pubsub/message"

type Command interface {
	// GetHeader returns the header of the command.
	GetHeader() message.Header
}
