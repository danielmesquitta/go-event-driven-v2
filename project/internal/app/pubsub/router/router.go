package router

import (
	"context"
	"tickets/internal/app/pubsub/handler/bookingcanceled"
	"tickets/internal/app/pubsub/handler/bookingconfirmed"
	"tickets/internal/app/pubsub/handler/bookingmade"
	"tickets/internal/app/pubsub/handler/refundticket"
	"tickets/internal/app/pubsub/handler/ticketprinted"
	"tickets/internal/app/pubsub/handler/ticketreceiptissued"
	"tickets/internal/app/pubsub/handler/ticketrefunded"
	pubSubMiddleware "tickets/internal/app/pubsub/middleware"
	"tickets/internal/pkg/log"
	"tickets/internal/provider/cmdbus"
	"tickets/internal/provider/eventbus"
	"time"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/message/router/middleware"
)

type Router struct {
	eventBus                                     eventbus.EventBus
	cmdBus                                       cmdbus.CommandBus
	appendCanceledBookingToTrackerHandler        *bookingcanceled.AppendToTracker
	appendConfirmedBookingToTrackerHandler       *bookingconfirmed.AppendToTracker
	issueReceiptHandler                          *bookingconfirmed.IssueReceipt
	createTicketHandler                          *bookingconfirmed.CreateTicket
	deleteTicketHandler                          *bookingcanceled.DeleteTicket
	printTicketHandler                           *bookingconfirmed.PrintTicket
	postTicketBookingToDeadNationHandler         *bookingmade.PostTicketBookingToDeadNation
	refundTicketHandler                          *refundticket.RefundTicket
	updateOpsBookingOnBookingConfirmedHandler    *bookingconfirmed.UpdateOpsBooking
	createOpsBookingHandler                      *bookingmade.CreateOpsBooking
	updateOpsBookingOnTicketRefundedHandler      *ticketrefunded.UpdateOpsBooking
	updateOpsBookingOnTicketPrintedHandler       *ticketprinted.UpdateOpsBooking
	updateOpsBookingOnTicketReceiptIssuedHandler *ticketreceiptissued.UpdateOpsBooking
}

func NewRouter(
	eventBus eventbus.EventBus,
	cmdBus cmdbus.CommandBus,
	appendCanceledBookingToTrackerHandler *bookingcanceled.AppendToTracker,
	appendConfirmedBookingToTrackerHandler *bookingconfirmed.AppendToTracker,
	issueReceiptHandler *bookingconfirmed.IssueReceipt,
	createTicketHandler *bookingconfirmed.CreateTicket,
	deleteTicketHandler *bookingcanceled.DeleteTicket,
	printTicketHandler *bookingconfirmed.PrintTicket,
	postTicketBookingToDeadNationHandler *bookingmade.PostTicketBookingToDeadNation,
	refundTicketHandler *refundticket.RefundTicket,
	updateOpsBookingOnBookingConfirmedHandler *bookingconfirmed.UpdateOpsBooking,
	createOpsBookingHandler *bookingmade.CreateOpsBooking,
	updateOpsBookingOnTicketRefundedHandler *ticketrefunded.UpdateOpsBooking,
	updateOpsBookingOnTicketPrintedHandler *ticketprinted.UpdateOpsBooking,
	updateOpsBookingOnTicketReceiptIssuedHandler *ticketreceiptissued.UpdateOpsBooking,
) *Router {
	return &Router{
		eventBus:                                     eventBus,
		cmdBus:                                       cmdBus,
		appendCanceledBookingToTrackerHandler:        appendCanceledBookingToTrackerHandler,
		appendConfirmedBookingToTrackerHandler:       appendConfirmedBookingToTrackerHandler,
		issueReceiptHandler:                          issueReceiptHandler,
		createTicketHandler:                          createTicketHandler,
		deleteTicketHandler:                          deleteTicketHandler,
		printTicketHandler:                           printTicketHandler,
		postTicketBookingToDeadNationHandler:         postTicketBookingToDeadNationHandler,
		refundTicketHandler:                          refundTicketHandler,
		updateOpsBookingOnBookingConfirmedHandler:    updateOpsBookingOnBookingConfirmedHandler,
		createOpsBookingHandler:                      createOpsBookingHandler,
		updateOpsBookingOnTicketRefundedHandler:      updateOpsBookingOnTicketRefundedHandler,
		updateOpsBookingOnTicketPrintedHandler:       updateOpsBookingOnTicketPrintedHandler,
		updateOpsBookingOnTicketReceiptIssuedHandler: updateOpsBookingOnTicketReceiptIssuedHandler,
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

	middlewares := []message.HandlerMiddleware{
		retry.Middleware,
		pubSubMiddleware.Logger,
		pubSubMiddleware.ErrorHandler,
		middleware.CorrelationID,
		pubSubMiddleware.CorrelationID,
		pubSubMiddleware.IdempotencyKey,
	}

	r.eventBus.AddMiddleware(middlewares...)
	r.cmdBus.AddMiddleware(middlewares...)

	// event.TicketBookingCanceled
	registerEventHandlers(
		r.eventBus,
		r.deleteTicketHandler.Handle,
		r.appendCanceledBookingToTrackerHandler.Handle,
	)

	// event.TicketBookingConfirmed
	registerEventHandlers(r.eventBus,
		r.appendConfirmedBookingToTrackerHandler.Handle,
		r.createTicketHandler.Handle,
		r.printTicketHandler.Handle,
		r.issueReceiptHandler.Handle,
		r.updateOpsBookingOnBookingConfirmedHandler.Handle,
	)

	// event.BookingMade
	registerEventHandlers(
		r.eventBus,
		r.postTicketBookingToDeadNationHandler.Handle,
		r.createOpsBookingHandler.Handle,
	)

	// event.TicketRefunded
	registerEventHandlers(
		r.eventBus,
		r.updateOpsBookingOnTicketRefundedHandler.Handle,
	)

	// event.TicketPrinted
	registerEventHandlers(
		r.eventBus,
		r.updateOpsBookingOnTicketPrintedHandler.Handle,
	)

	// event.TicketReceiptIssued
	registerEventHandlers(
		r.eventBus,
		r.updateOpsBookingOnTicketReceiptIssuedHandler.Handle,
	)

	// cmd.RefundTicket
	registerCommandHandlers(r.cmdBus,
		r.refundTicketHandler.Handle,
	)

	go func() {
		err := r.eventBus.Run(ctx)
		if err != nil {
			panic(err)
		}
	}()

	go func() {
		err := r.cmdBus.Run(ctx)
		if err != nil {
			panic(err)
		}
	}()

	return nil
}

func (r *Router) Running() chan struct{} {
	ch := make(chan struct{})
	go func() {
		<-r.eventBus.Running()
		<-r.cmdBus.Running()
		close(ch)
	}()
	return ch
}

func registerEventHandlers[T any](
	eventBus eventbus.EventBus,
	handlerFuncs ...eventbus.HandlerFunc[T],
) {
	for _, handlerFunc := range handlerFuncs {
		eventbus.AddHandler(eventBus, handlerFunc)
	}
}

func registerCommandHandlers[T any](
	cmdBus cmdbus.CommandBus,
	handlerFuncs ...cmdbus.HandlerFunc[T],
) {
	for _, handlerFunc := range handlerFuncs {
		cmdbus.AddHandler(cmdBus, handlerFunc)
	}
}
