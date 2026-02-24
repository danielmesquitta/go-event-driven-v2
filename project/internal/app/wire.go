//go:build wireinject

package app

import (
	"github.com/goforj/wire"

	"tickets/internal/app/http/handler"
	httpRouter "tickets/internal/app/http/router"
	"tickets/internal/app/pubsub/handler/bookingcanceled"
	"tickets/internal/app/pubsub/handler/bookingconfirmed"
	pubSubRouter "tickets/internal/app/pubsub/router"
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
		gateway.NewGateway,
		db.NewDB,
		redisstream.NewEventBus,

		pg.NewTicketRepo,
		pg.NewShowRepo,
		wire.Bind(new(repo.TicketRepo), new(*pg.TicketRepo)),
		wire.Bind(new(repo.ShowRepo), new(*pg.ShowRepo)),

		spreadsheetapi.NewClient,
		receiptsvc.NewClient,
		filestorageapi.NewClient,
		wire.Bind(new(spreadsheet.API), new(*spreadsheetapi.Client)),
		wire.Bind(new(receipt.Service), new(*receiptsvc.Client)),
		wire.Bind(new(filestorage.Storage), new(*filestorageapi.Client)),
		wire.Bind(new(eventbus.EventBus), new(*redisstream.EventBus)),

		tracker.NewAppendCanceledBooking,
		tracker.NewAppendConfirmedBooking,
		ticket.NewIssueReceipt,
		ticket.NewCreate,
		ticket.NewDelete,
		ticket.NewPrint,
		ticket.NewPostStatus,
		ticket.NewList,
		show.NewCreate,

		bookingcanceled.NewAppendToTracker,
		bookingcanceled.NewDeleteTicket,
		bookingconfirmed.NewAppendToTracker,
		bookingconfirmed.NewIssueReceipt,
		bookingconfirmed.NewCreateTicket,
		bookingconfirmed.NewPrintTicket,

		pubSubRouter.NewRouter,
		httpRouter.New,
		handler.NewShowHandler,
		handler.NewTicketHandler,
		handler.NewHealthHandler,

		wire.Struct(new(Service), "*"),
	)
	return Service{}
}
