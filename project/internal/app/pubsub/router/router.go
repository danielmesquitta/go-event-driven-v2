package router

import (
	"context"
	"log/slog"
	"tickets/internal/app/pubsub/event"
	"tickets/internal/provider/pubsub"
	"tickets/internal/provider/receipt"
	"tickets/internal/provider/spreadsheet"

	"github.com/ThreeDotsLabs/watermill/message"
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
	r.pubsub.AddMiddleware(middleware.CorrelationID, LogMiddleware)

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

func LogMiddleware(next message.HandlerFunc) message.HandlerFunc {
	return func(msg *message.Message) ([]*message.Message, error) {
		logger := slog.With(
			"message_id", msg.UUID,
			"payload", string(msg.Payload),
			"metadata", msg.Metadata,
			"handler", message.HandlerNameFromCtx(msg.Context()),
		)

		logger.Info("Handling a message")

		return next(msg)
	}
}
