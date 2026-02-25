//go:build wireinject

package app

import (
	"github.com/goforj/wire"

	"tickets/internal/app/http/handler"
	httpRouter "tickets/internal/app/http/router"
	"tickets/internal/app/pubsub/handler/bookingcanceled"
	"tickets/internal/app/pubsub/handler/bookingconfirmed"
	pubSubRouter "tickets/internal/app/pubsub/router"
	"tickets/internal/domain/usecase/booking"
	"tickets/internal/domain/usecase/show"
	"tickets/internal/domain/usecase/ticket"
	"tickets/internal/domain/usecase/tracker"
	"tickets/internal/provider/db"
	"tickets/internal/provider/eventbus"
	"tickets/internal/provider/eventbus/redisstream"
	"tickets/internal/provider/filestorage"
	"tickets/internal/provider/filestorage/filestorageapi"
	"tickets/internal/provider/gateway"
	"tickets/internal/provider/receipt"
	"tickets/internal/provider/receipt/receiptsvc"
	"tickets/internal/provider/repo"
	"tickets/internal/provider/repo/pg"
	"tickets/internal/provider/spreadsheet"
	"tickets/internal/provider/spreadsheet/spreadsheetapi"
)

func New() Service {
	wire.Build(
		// Repos
		pg.NewTicketRepo,
		wire.Bind(new(repo.TicketRepo), new(*pg.TicketRepo)),
		pg.NewShowRepo,
		wire.Bind(new(repo.ShowRepo), new(*pg.ShowRepo)),
		pg.NewBookingRepo,
		wire.Bind(new(repo.BookingRepo), new(*pg.BookingRepo)),

		// Providers
		gateway.NewGateway,
		db.NewDB,
		filestorageapi.NewClient,
		wire.Bind(new(spreadsheet.API), new(*spreadsheetapi.Client)),
		receiptsvc.NewClient,
		wire.Bind(new(receipt.Service), new(*receiptsvc.Client)),
		spreadsheetapi.NewClient,
		wire.Bind(new(filestorage.Storage), new(*filestorageapi.Client)),
		redisstream.NewEventBus,
		wire.Bind(new(eventbus.EventBus), new(*redisstream.EventBus)),

		// Usecases
		tracker.NewAppendCanceledBooking,
		tracker.NewAppendConfirmedBooking,
		ticket.NewIssueReceipt,
		ticket.NewCreate,
		ticket.NewDelete,
		ticket.NewPrint,
		ticket.NewPostStatus,
		ticket.NewList,
		show.NewCreate,
		booking.NewCreate,

		// PubSub handlers
		bookingcanceled.NewAppendToTracker,
		bookingcanceled.NewDeleteTicket,
		bookingconfirmed.NewAppendToTracker,
		bookingconfirmed.NewIssueReceipt,
		bookingconfirmed.NewCreateTicket,
		bookingconfirmed.NewPrintTicket,

		// HTTP handlers
		pubSubRouter.NewRouter,
		httpRouter.New,
		handler.NewShowHandler,
		handler.NewTicketHandler,
		handler.NewHealthHandler,
		handler.NewBookingHandler,

		wire.Struct(new(Service), "*"),
	)
	return Service{}
}
