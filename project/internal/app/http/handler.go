package http

import "tickets/internal/provider/pubsub"

type Handler struct {
	pubsub pubsub.PubSub
}
