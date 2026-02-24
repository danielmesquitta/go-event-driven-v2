package app

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"

	"github.com/labstack/echo/v4"
	"golang.org/x/sync/errgroup"

	"tickets/internal/app/http/handler"
	httpRouter "tickets/internal/app/http/router"
	"tickets/internal/app/pubsub/handler/bookingcanceled"
	"tickets/internal/app/pubsub/handler/bookingconfirmed"
	pubSubRouter "tickets/internal/app/pubsub/router"
	"tickets/internal/domain/usecase/show"
	"tickets/internal/domain/usecase/ticket"
	"tickets/internal/domain/usecase/tracker"
	"tickets/internal/provider/db"
	"tickets/internal/provider/eventbus/redisstream"
	"tickets/internal/provider/filestorage"
	"tickets/internal/provider/filestorage/filestorageapi"
	"tickets/internal/provider/gateway"
	"tickets/internal/provider/receipt"
	"tickets/internal/provider/receipt/receiptsvc"
	"tickets/internal/provider/repo/pg"
	"tickets/internal/provider/spreadsheet"
	"tickets/internal/provider/spreadsheet/spreadsheetapi"
)

type Service struct {
	db              *db.DB
	httpRouter      *echo.Echo
	pubSubRouter    *pubSubRouter.Router
	spreadsheetsAPI spreadsheet.API
	receiptsService receipt.Service
	fileStorage     filestorage.Storage
}

func New() Service {
	gateway := gateway.NewGateway()

	db := db.NewDB()
	ticketRepo := pg.NewTicketRepo(db)
	showRepo := pg.NewShowRepo(db)

	spreadsheetsAPI := spreadsheetapi.NewClient(gateway)
	receiptsService := receiptsvc.NewClient(gateway)
	fileStorageClient := filestorageapi.NewClient(gateway)

	eventBus := redisstream.NewEventBus()

	appendCanceledBookingToTrackerUseCase := tracker.NewAppendCanceledBooking(spreadsheetsAPI)
	appendCanceledBookingToTrackerHandler := bookingcanceled.NewAppendToTracker(appendCanceledBookingToTrackerUseCase)

	appendConfirmedBookingToTrackerUseCase := tracker.NewAppendConfirmedBooking(spreadsheetsAPI)
	appendConfirmedBookingToTrackerHandler := bookingconfirmed.NewAppendToTracker(appendConfirmedBookingToTrackerUseCase)

	issueReceiptUseCase := ticket.NewIssueReceipt(receiptsService)
	issueReceiptHandler := bookingconfirmed.NewIssueReceipt(issueReceiptUseCase)

	createTicketUseCase := ticket.NewCreate(ticketRepo)
	createTicketHandler := bookingconfirmed.NewCreateTicket(createTicketUseCase)

	deleteTicketUseCase := ticket.NewDelete(ticketRepo)
	deleteTicketHandler := bookingcanceled.NewDeleteTicket(deleteTicketUseCase)

	printTicketUseCase := ticket.NewPrint(fileStorageClient, eventBus)
	printTicketHandler := bookingconfirmed.NewPrintTicket(printTicketUseCase)

	pubSubRouter := pubSubRouter.NewRouter(
		eventBus,
		appendCanceledBookingToTrackerHandler,
		appendConfirmedBookingToTrackerHandler,
		issueReceiptHandler,
		createTicketHandler,
		deleteTicketHandler,
		printTicketHandler,
	)

	postTicketStatusUseCase := ticket.NewPostStatus(eventBus)
	listTicketsUseCase := ticket.NewList(ticketRepo)

	createShowUseCase := show.NewCreate(showRepo)
	showHandler := handler.NewShowHandler(createShowUseCase)

	ticketHandler := handler.NewTicketHandler(postTicketStatusUseCase, listTicketsUseCase)

	healthHandler := handler.NewHealthHandler()

	httpRouter := httpRouter.New(
		showHandler,
		ticketHandler,
		healthHandler,
	)

	svc := Service{
		db:              db,
		httpRouter:      httpRouter,
		pubSubRouter:    pubSubRouter,
		spreadsheetsAPI: spreadsheetsAPI,
		receiptsService: receiptsService,
		fileStorage:     fileStorageClient,
	}

	return svc
}

func (s Service) Run(ctx context.Context) error {
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()

	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		return s.db.InitializeSchema(ctx)
	})

	g.Go(func() error {
		return s.pubSubRouter.Run(ctx)
	})

	g.Go(func() error {
		<-s.pubSubRouter.Running()

		err := s.httpRouter.Start(":8080")
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			return err
		}

		return nil
	})

	g.Go(func() error {
		<-ctx.Done()
		return s.httpRouter.Shutdown(ctx)
	})

	return g.Wait()
}
