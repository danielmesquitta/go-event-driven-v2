package router

import (
	"context"
	"tickets/internal/app/pubsub/event"
	pubSubMiddleware "tickets/internal/app/pubsub/middleware"
	"tickets/internal/provider/pubsub"
	"tickets/internal/provider/receipt"
	"tickets/internal/provider/spreadsheet"

	"github.com/ThreeDotsLabs/watermill/message/router/middleware"
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
	r.pubsub.AddMiddleware(
		middleware.CorrelationID,
		pubSubMiddleware.CorrelationID,
		pubSubMiddleware.Logger,
	)

	r.pubsub.AddConsumerHandler(
		event.TopicTicketBookingConfirmed,
		r.handleIssueReceipt,
	)

	r.pubsub.AddConsumerHandler(
		event.TopicTicketBookingConfirmed,
		r.handleAppendConfirmedBookingToTracker,
	)

	r.pubsub.AddConsumerHandler(
		event.TopicTicketBookingCanceled,
		r.handleAppendCanceledBookingToTracker,
	)

	return r.pubsub.Run(ctx)
}

func (r *Router) Running() chan struct{} {
	return r.pubsub.Running()
}
