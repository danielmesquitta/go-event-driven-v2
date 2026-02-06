package pubsub

import (
	"context"
	"tickets/internal/provider/pubsub"
	"tickets/internal/provider/receipt"
	"tickets/internal/provider/spreadsheet"
)

type Router struct {
	pubsub pubsub.PubSub

	spreadsheetAPI  spreadsheet.API
	receiptsService receipt.Service
}

func NewRouter(
	pubsub pubsub.PubSub,
) *Router {
	return &Router{
		pubsub: pubsub,
	}
}

func (r *Router) Run(
	ctx context.Context,
) error {
	r.pubsub.AddConsumerHandler(
		"HandleIssueReceipt",
		pubsub.TopicIssueReceipt,
		r.handleIssueReceipt,
	)

	r.pubsub.AddConsumerHandler(
		"HandleAppendToTracker",
		pubsub.TopicAppendToTracker,
		r.handleAppendToTracker,
	)

	return r.pubsub.Register(ctx)
}

func (r *Router) Running() chan struct{} {
	return r.pubsub.Running()
}
