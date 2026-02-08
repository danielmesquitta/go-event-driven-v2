package app

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"

	"github.com/ThreeDotsLabs/go-event-driven/v2/common/clients"
	"github.com/ThreeDotsLabs/go-event-driven/v2/common/log"
	"github.com/labstack/echo/v4"
	"golang.org/x/sync/errgroup"

	appHTTP "tickets/internal/app/http"
	"tickets/internal/app/pubsub/handler"
	"tickets/internal/app/pubsub/router"
	"tickets/internal/provider/eventbus/redisstream"
	"tickets/internal/provider/receipt"
	"tickets/internal/provider/receipt/receiptsvc"
	"tickets/internal/provider/spreadsheet"
	"tickets/internal/provider/spreadsheet/spreadsheetapi"
)

type Service struct {
	echoRouter      *echo.Echo
	router          *router.Router
	spreadsheetsAPI spreadsheet.API
	receiptsService receipt.Service
}

func New() Service {
	log.Init(slog.LevelInfo)

	apiClients, err := clients.NewClients(
		os.Getenv("GATEWAY_ADDR"),
		func(ctx context.Context, req *http.Request) error {
			req.Header.Set("Correlation-ID", log.CorrelationIDFromContext(ctx))
			return nil
		},
	)
	if err != nil {
		panic(err)
	}

	spreadsheetsAPI := spreadsheetapi.NewClient(apiClients)
	receiptsService := receiptsvc.NewClient(apiClients)

	eventBus := redisstream.NewEventBus()

	appendCanceledBookingToTrackerHandler := handler.NewAppendCanceledBookingToTracker(spreadsheetsAPI)
	appendConfirmedBookingToTrackerHandler := handler.NewAppendConfirmedBookingToTracker(spreadsheetsAPI)
	issueReceiptHandler := handler.NewIssueReceipt(receiptsService)

	router := router.NewRouter(
		eventBus,
		appendCanceledBookingToTrackerHandler,
		appendConfirmedBookingToTrackerHandler,
		issueReceiptHandler,
	)

	echoRouter := appHTTP.NewHttpRouter(eventBus)

	svc := Service{
		echoRouter:      echoRouter,
		router:          router,
		spreadsheetsAPI: spreadsheetsAPI,
		receiptsService: receiptsService,
	}

	return svc
}

func (s Service) Run(ctx context.Context) error {
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()

	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		return s.router.Run(ctx)
	})

	g.Go(func() error {
		<-s.router.Running()

		err := s.echoRouter.Start(":8080")
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			return err
		}

		return nil
	})

	g.Go(func() error {
		<-ctx.Done()
		return s.echoRouter.Shutdown(ctx)
	})

	return g.Wait()
}
