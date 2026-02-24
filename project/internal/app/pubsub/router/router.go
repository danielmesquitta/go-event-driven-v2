package router

import (
	"context"
	"tickets/internal/app/pubsub/event"
	"tickets/internal/app/pubsub/handler/bookingcanceled"
	"tickets/internal/app/pubsub/handler/bookingconfirmed"
	pubSubMiddleware "tickets/internal/app/pubsub/middleware"
	"tickets/internal/pkg/log"
	"tickets/internal/provider/eventbus"
	"time"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message/router/middleware"
)

type Router struct {
	eventBus                              eventbus.EventBus
	appendCanceledBookingToTrackerHandler *bookingcanceled.AppendToTracker
	appendConfirmedBookingToTrackerHandler *bookingconfirmed.AppendToTracker
	issueReceiptHandler                   *bookingconfirmed.IssueReceipt
	createTicketHandler                   *bookingconfirmed.CreateTicket
	deleteTicketHandler                   *bookingcanceled.DeleteTicket
	printTicketHandler                    *bookingconfirmed.PrintTicket
}

func NewRouter(
	eventBus eventbus.EventBus,
	appendCanceledBookingToTrackerHandler *bookingcanceled.AppendToTracker,
	appendConfirmedBookingToTrackerHandler *bookingconfirmed.AppendToTracker,
	issueReceiptHandler *bookingconfirmed.IssueReceipt,
	createTicketHandler *bookingconfirmed.CreateTicket,
	deleteTicketHandler *bookingcanceled.DeleteTicket,
	printTicketHandler *bookingconfirmed.PrintTicket,
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

	r.registerTicketBookingCanceledHandlers([]eventbus.HandlerFunc[event.TicketBookingCanceled]{
		r.deleteTicketHandler.Handle,
		r.appendCanceledBookingToTrackerHandler.Handle,
	})

	r.registerTicketBookingConfirmedHandlers([]eventbus.HandlerFunc[event.TicketBookingConfirmed]{
		r.appendConfirmedBookingToTrackerHandler.Handle,
		r.createTicketHandler.Handle,
		r.printTicketHandler.Handle,
		r.issueReceiptHandler.Handle,
	})

	return r.eventBus.Run(ctx)
}

func (r *Router) Running() chan struct{} {
	return r.eventBus.Running()
}

func (r *Router) registerTicketBookingCanceledHandlers(
	handlerFuncs []eventbus.HandlerFunc[event.TicketBookingCanceled],
) {
	registerHandlers(r.eventBus, handlerFuncs)
}

func (r *Router) registerTicketBookingConfirmedHandlers(
	handlerFuncs []eventbus.HandlerFunc[event.TicketBookingConfirmed],
) {
	registerHandlers(r.eventBus, handlerFuncs)
}

func registerHandlers[T any](
	eventBus eventbus.EventBus,
	handlerFuncs []eventbus.HandlerFunc[T],
) {
	for _, handlerFunc := range handlerFuncs {
		eventbus.AddHandler(eventBus, handlerFunc)
	}
}
