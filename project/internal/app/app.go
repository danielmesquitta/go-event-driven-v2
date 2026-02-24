package app

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"

	"github.com/ThreeDotsLabs/go-event-driven/v2/common/clients"
	"github.com/labstack/echo/v4"
	"golang.org/x/sync/errgroup"

	httpRouter "tickets/internal/app/http/router"
	"tickets/internal/app/pubsub/handler"
	pubSubRouter "tickets/internal/app/pubsub/router"
	"tickets/internal/domain/usecase/ticket"
	"tickets/internal/domain/usecase/tracker"
	"tickets/internal/pkg/ctxval"
	"tickets/internal/provider/db"
	"tickets/internal/provider/eventbus/redisstream"
	"tickets/internal/provider/filestorage"
	"tickets/internal/provider/filestorage/filestorageapi"
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
	apiClients, err := clients.NewClients(
		os.Getenv("GATEWAY_ADDR"),
		func(ctx context.Context, req *http.Request) error {
			req.Header.Set("Correlation-ID", ctxval.GetCorrelationID(ctx))
			return nil
		},
	)
	if err != nil {
		panic(err)
	}

	db := db.NewDB()
	ticketRepo := pg.NewTicketRepo(db)

	spreadsheetsAPI := spreadsheetapi.NewClient(apiClients)
	receiptsService := receiptsvc.NewClient(apiClients)
	fileStorageClient := filestorageapi.NewClient(apiClients)

	eventBus := redisstream.NewEventBus()

	appendCanceledBookingToTrackerUseCase := tracker.NewAppendCanceledBooking(spreadsheetsAPI)
	appendCanceledBookingToTrackerHandler := handler.NewAppendCanceledBookingToTracker(appendCanceledBookingToTrackerUseCase)

	appendConfirmedBookingToTrackerUseCase := tracker.NewAppendConfirmedBooking(spreadsheetsAPI)
	appendConfirmedBookingToTrackerHandler := handler.NewAppendConfirmedBookingToTracker(appendConfirmedBookingToTrackerUseCase)

	issueReceiptUseCase := ticket.NewIssueReceipt(receiptsService)
	issueReceiptHandler := handler.NewIssueReceipt(issueReceiptUseCase)

	createTicketUseCase := ticket.NewCreate(ticketRepo)
	createTicketHandler := handler.NewCreateTicket(createTicketUseCase)

	deleteTicketUseCase := ticket.NewDelete(ticketRepo)
	deleteTicketHandler := handler.NewDeleteTicket(deleteTicketUseCase)

	printTicketUseCase := ticket.NewPrint(fileStorageClient, eventBus)
	printTicketHandler := handler.NewPrintTicket(printTicketUseCase)

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

	httpRouter := httpRouter.New(
		postTicketStatusUseCase,
		listTicketsUseCase,
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
