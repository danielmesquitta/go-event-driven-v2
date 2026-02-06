package pubsub

import (
	"context"
	"tickets/internal/app/pubsub/event"
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
	spreadsheetAPI spreadsheet.API,
	receiptsService receipt.Service,
) *Router {
	return &Router{
		pubsub:          pubsub,
		spreadsheetAPI:  spreadsheetAPI,
		receiptsService: receiptsService,
	}
}

func (r *Router) Run(
	ctx context.Context,
) error {
	r.pubsub.AddConsumerHandler(
		event.TopicTicketBookingConfirmed,
		r.handleIssueReceipt,
	)

	r.pubsub.AddConsumerHandler(
		event.TopicTicketBookingConfirmed,
		r.handleAppendToTracker,
	)

	return r.pubsub.Register(ctx)
}

func (r *Router) Running() chan struct{} {
	return r.pubsub.Running()
}
