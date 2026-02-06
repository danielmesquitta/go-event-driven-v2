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
	"tickets/internal/provider/pubsub"
	"tickets/internal/provider/pubsub/redisstream"
	"tickets/internal/provider/receipt"
	"tickets/internal/provider/receipt/receiptsvc"
	"tickets/internal/provider/spreadsheet"
	"tickets/internal/provider/spreadsheet/spreadsheetapi"
)

type Service struct {
	echoRouter      *echo.Echo
	pubsub          pubsub.PubSub
	spreadsheetsAPI spreadsheet.API
	receiptsService receipt.Service
}

func New() Service {
	log.Init(slog.LevelInfo)

	apiClients, err := clients.NewClients(os.Getenv("GATEWAY_ADDR"), nil)
	if err != nil {
		panic(err)
	}

	spreadsheetsAPI := spreadsheetapi.NewClient(apiClients)
	receiptsService := receiptsvc.NewClient(apiClients)

	pubsub := redisstream.NewPubSub(spreadsheetsAPI, receiptsService)

	echoRouter := appHTTP.NewHttpRouter(pubsub.Publisher)

	svc := Service{
		echoRouter:      echoRouter,
		pubsub:          pubsub,
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
		return s.pubsub.Register(ctx)
	})

	g.Go(func() error {
		<-s.pubsub.Running()

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
