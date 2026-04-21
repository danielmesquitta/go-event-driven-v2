//go:build wireinject

package app

import (
	"github.com/goforj/wire"

	"tickets/internal/app/http/handler"
	httpRouter "tickets/internal/app/http/router"
	"tickets/internal/app/pubsub/handler/bookingcanceled"
	"tickets/internal/app/pubsub/handler/bookingconfirmed"
	"tickets/internal/app/pubsub/handler/bookingmade"
	"tickets/internal/app/pubsub/handler/refundticket"
	"tickets/internal/app/pubsub/handler/ticketprinted"
	"tickets/internal/app/pubsub/handler/ticketreceiptissued"
	"tickets/internal/app/pubsub/handler/ticketrefunded"
	pubSubRouter "tickets/internal/app/pubsub/router"
	"tickets/internal/domain/usecase/booking"
	"tickets/internal/domain/usecase/opsbooking"
	"tickets/internal/domain/usecase/show"
	"tickets/internal/domain/usecase/ticket"
	"tickets/internal/domain/usecase/tracker"
	"tickets/internal/pkg/tx"
	"tickets/internal/provider/cmdbus"
	cmdBusRedisStream "tickets/internal/provider/cmdbus/redisstream"
	"tickets/internal/provider/db"
	"tickets/internal/provider/eventbus"
	"tickets/internal/provider/eventbus/redisstream"
	"tickets/internal/provider/filestorage"
	"tickets/internal/provider/filestorage/filestorageapi"
	"tickets/internal/provider/gateway"
	"tickets/internal/provider/outbox"
	pgOutbox "tickets/internal/provider/outbox/pg"
	"tickets/internal/provider/payment"
	"tickets/internal/provider/payment/paymentsvc"
	"tickets/internal/provider/receipt"
	"tickets/internal/provider/receipt/receiptsvc"
	"tickets/internal/provider/repo"
	pgRepo "tickets/internal/provider/repo/pg"
	"tickets/internal/provider/showapi"
	"tickets/internal/provider/showapi/deadnation"
	"tickets/internal/provider/spreadsheet"
	"tickets/internal/provider/spreadsheet/spreadsheetapi"
)

func New() Service {
	wire.Build(
		// Repos
		pgRepo.NewTicketRepo,
		wire.Bind(new(repo.TicketRepo), new(*pgRepo.TicketRepo)),
		pgRepo.NewShowRepo,
		wire.Bind(new(repo.ShowRepo), new(*pgRepo.ShowRepo)),
		pgRepo.NewBookingRepo,
		wire.Bind(new(repo.BookingRepo), new(*pgRepo.BookingRepo)),
		pgRepo.NewOpsBooking,
		wire.Bind(new(repo.OpsBookingRepo), new(*pgRepo.OpsBooking)),

		// Packages
		tx.NewSqlxTransaction,
		wire.Bind(new(tx.Transaction), new(*tx.SqlxTransaction)),

		// Providers
		gateway.NewGateway,
		db.NewSQLXDB,
		db.NewDB,
		filestorageapi.NewClient,
		wire.Bind(new(spreadsheet.API), new(*spreadsheetapi.Client)),
		receiptsvc.NewClient,
		wire.Bind(new(receipt.Service), new(*receiptsvc.Client)),
		spreadsheetapi.NewClient,
		wire.Bind(new(filestorage.Storage), new(*filestorageapi.Client)),
		redisstream.NewEventBus,
		wire.Bind(new(eventbus.EventBus), new(*redisstream.EventBus)),
		pgOutbox.New,
		wire.Bind(new(outbox.Outbox), new(*pgOutbox.Outbox)),
		deadnation.New,
		wire.Bind(new(showapi.ShowAPI), new(*deadnation.DeadNationAPI)),
		cmdBusRedisStream.NewCommandBus,
		wire.Bind(new(cmdbus.CommandBus), new(*cmdBusRedisStream.CommandBus)),
		paymentsvc.NewClient,
		wire.Bind(new(payment.Service), new(*paymentsvc.Client)),

		// Usecases
		tracker.NewAppendCanceledBooking,
		tracker.NewAppendConfirmedBooking,
		ticket.NewIssueReceipt,
		ticket.NewCreate,
		ticket.NewDelete,
		ticket.NewPrint,
		ticket.NewPostStatus,
		ticket.NewList,
		ticket.NewSendRefundCommand,
		ticket.NewRefund,
		show.NewCreate,
		booking.NewCreate,
		booking.NewPostTicketBookingToDeadNation,
		opsbooking.NewCreate,

		// PubSub handlers
		bookingcanceled.NewAppendToTracker,
		bookingcanceled.NewDeleteTicket,
		bookingconfirmed.NewAppendToTracker,
		bookingconfirmed.NewIssueReceipt,
		bookingconfirmed.NewCreateTicket,
		bookingconfirmed.NewPrintTicket,
		bookingconfirmed.NewUpdateOpsBooking,
		bookingmade.NewPostTicketBookingToDeadNation,
		bookingmade.NewCreateOpsBooking,
		refundticket.NewRefundTicket,
		ticketprinted.NewUpdateOpsBooking,
		ticketreceiptissued.NewUpdateOpsBooking,
		ticketrefunded.NewUpdateOpsBooking,

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
