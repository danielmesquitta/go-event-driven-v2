package router

import (
	"context"
	"tickets/internal/app/pubsub/event"
	"tickets/internal/app/pubsub/handler"
	pubSubMiddleware "tickets/internal/app/pubsub/middleware"
	"tickets/internal/pkg/log"
	"tickets/internal/provider/pubsub"
	"time"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message/router/middleware"
)

type Router struct {
	pubsub pubsub.PubSub

	appendCanceledBookingToTrackerHandler  *handler.AppendCanceledBookingToTracker
	appendConfirmedBookingToTrackerHandler *handler.AppendConfirmedBookingToTracker
	issueReceiptHandler                    *handler.IssueReceipt
}

func NewRouter(
	pubsub pubsub.PubSub,
	appendCanceledBookingToTrackerHandler *handler.AppendCanceledBookingToTracker,
	appendConfirmedBookingToTrackerHandler *handler.AppendConfirmedBookingToTracker,
	issueReceiptHandler *handler.IssueReceipt,
) *Router {
	return &Router{
		pubsub:                                 pubsub,
		appendCanceledBookingToTrackerHandler:  appendCanceledBookingToTrackerHandler,
		appendConfirmedBookingToTrackerHandler: appendConfirmedBookingToTrackerHandler,
		issueReceiptHandler:                    issueReceiptHandler,
	}
}

func (r *Router) Run(
	ctx context.Context,
) error {
	retry := middleware.Retry{
		MaxRetries:      10,
		InitialInterval: time.Millisecond * 100,
		MaxInterval:     time.Second,
		Multiplier:      2,
		Logger:          watermill.NewSlogLogger(log.FromContext(ctx)),
	}

	r.pubsub.AddMiddleware(
		middleware.CorrelationID,
		retry.Middleware,
		pubSubMiddleware.CorrelationID,
		pubSubMiddleware.Logger,
		pubSubMiddleware.ErrorHandler,
	)

	r.pubsub.AddConsumerHandler(
		event.TopicTicketBookingConfirmed,
		r.issueReceiptHandler.Handle,
		r.appendConfirmedBookingToTrackerHandler.Handle,
	)

	r.pubsub.AddConsumerHandler(
		event.TopicTicketBookingCanceled,
		r.appendCanceledBookingToTrackerHandler.Handle,
	)

	return r.pubsub.Run(ctx)
}

func (r *Router) Running() chan struct{} {
	return r.pubsub.Running()
}
