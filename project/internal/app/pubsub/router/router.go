package router

import (
	"context"
	"tickets/internal/app/pubsub/handler/bookingcanceled"
	"tickets/internal/app/pubsub/handler/bookingconfirmed"
	"tickets/internal/app/pubsub/handler/bookingmade"
	pubSubMiddleware "tickets/internal/app/pubsub/middleware"
	"tickets/internal/pkg/log"
	"tickets/internal/provider/eventbus"
	"time"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message/router/middleware"
)

type Router struct {
	eventBus                               eventbus.EventBus
	appendCanceledBookingToTrackerHandler  *bookingcanceled.AppendToTracker
	appendConfirmedBookingToTrackerHandler *bookingconfirmed.AppendToTracker
	issueReceiptHandler                    *bookingconfirmed.IssueReceipt
	createTicketHandler                    *bookingconfirmed.CreateTicket
	deleteTicketHandler                    *bookingcanceled.DeleteTicket
	printTicketHandler                     *bookingconfirmed.PrintTicket
	mapShowIdToDeadNationEventIdHandler    *bookingmade.PostTicketBookingToDeadNation
}

func NewRouter(
	eventBus eventbus.EventBus,
	appendCanceledBookingToTrackerHandler *bookingcanceled.AppendToTracker,
	appendConfirmedBookingToTrackerHandler *bookingconfirmed.AppendToTracker,
	issueReceiptHandler *bookingconfirmed.IssueReceipt,
	createTicketHandler *bookingconfirmed.CreateTicket,
	deleteTicketHandler *bookingcanceled.DeleteTicket,
	printTicketHandler *bookingconfirmed.PrintTicket,
	mapShowIdToDeadNationEventIdHandler *bookingmade.PostTicketBookingToDeadNation,
) *Router {
	return &Router{
		eventBus:                               eventBus,
		appendCanceledBookingToTrackerHandler:  appendCanceledBookingToTrackerHandler,
		appendConfirmedBookingToTrackerHandler: appendConfirmedBookingToTrackerHandler,
		issueReceiptHandler:                    issueReceiptHandler,
		createTicketHandler:                    createTicketHandler,
		deleteTicketHandler:                    deleteTicketHandler,
		printTicketHandler:                     printTicketHandler,
		mapShowIdToDeadNationEventIdHandler:    mapShowIdToDeadNationEventIdHandler,
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

	// event.TicketBookingCanceled
	registerHandlers(
		r.eventBus,
		r.deleteTicketHandler.Handle,
		r.appendCanceledBookingToTrackerHandler.Handle,
	)

	// event.TicketBookingConfirmed
	registerHandlers(r.eventBus,
		r.appendConfirmedBookingToTrackerHandler.Handle,
		r.createTicketHandler.Handle,
		r.printTicketHandler.Handle,
		r.issueReceiptHandler.Handle,
	)

	// event.BookingMade
	registerHandlers(
		r.eventBus,
		r.mapShowIdToDeadNationEventIdHandler.Handle,
	)

	return r.eventBus.Run(ctx)
}

func (r *Router) Running() chan struct{} {
	return r.eventBus.Running()
}

func registerHandlers[T any](
	eventBus eventbus.EventBus,
	handlerFuncs ...eventbus.HandlerFunc[T],
) {
	for _, handlerFunc := range handlerFuncs {
		eventbus.AddHandler(eventBus, handlerFunc)
	}
}
