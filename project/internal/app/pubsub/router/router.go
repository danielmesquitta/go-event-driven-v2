package router

import (
	"context"
	"tickets/internal/app/pubsub/handler"
	pubSubMiddleware "tickets/internal/app/pubsub/middleware"
	"tickets/internal/pkg/log"
	"tickets/internal/provider/eventbus"
	"time"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message/router/middleware"
)

type Router struct {
	eventBus                               eventbus.EventBus
	appendCanceledBookingToTrackerHandler  *handler.AppendCanceledBookingToTracker
	appendConfirmedBookingToTrackerHandler *handler.AppendConfirmedBookingToTracker
	issueReceiptHandler                    *handler.IssueReceipt
	createTicketHandler                    *handler.CreateTicket
	deleteTicketHandler                    *handler.DeleteTicket
	printTicketHandler                     *handler.PrintTicket
}

func NewRouter(
	eventBus eventbus.EventBus,
	appendCanceledBookingToTrackerHandler *handler.AppendCanceledBookingToTracker,
	appendConfirmedBookingToTrackerHandler *handler.AppendConfirmedBookingToTracker,
	issueReceiptHandler *handler.IssueReceipt,
	createTicketHandler *handler.CreateTicket,
	deleteTicketHandler *handler.DeleteTicket,
	printTicketHandler *handler.PrintTicket,
) *Router {
	return &Router{
		eventBus:                               eventBus,
		appendCanceledBookingToTrackerHandler:  appendCanceledBookingToTrackerHandler,
		appendConfirmedBookingToTrackerHandler: appendConfirmedBookingToTrackerHandler,
		issueReceiptHandler:                    issueReceiptHandler,
		createTicketHandler:                    createTicketHandler,
		deleteTicketHandler:                    deleteTicketHandler,
		printTicketHandler:                     printTicketHandler,
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
		Logger:          watermill.NewSlogLogger(log.New(ctx)),
	}

	r.eventBus.AddMiddleware(
		retry.Middleware,
		pubSubMiddleware.Logger,
		pubSubMiddleware.ErrorHandler,
		middleware.CorrelationID,
		pubSubMiddleware.CorrelationID,
		pubSubMiddleware.IdempotencyKey,
	)

	// TicketBookingCanceled
	eventbus.AddHandler(r.eventBus, r.issueReceiptHandler.Handle)
	eventbus.AddHandler(r.eventBus, r.deleteTicketHandler.Handle)

	// TicketBookingConfirmed
	eventbus.AddHandler(r.eventBus, r.appendCanceledBookingToTrackerHandler.Handle)
	eventbus.AddHandler(r.eventBus, r.appendConfirmedBookingToTrackerHandler.Handle)
	eventbus.AddHandler(r.eventBus, r.createTicketHandler.Handle)
	eventbus.AddHandler(r.eventBus, r.printTicketHandler.Handle)

	return r.eventBus.Run(ctx)
}

func (r *Router) Running() chan struct{} {
	return r.eventBus.Running()
}
