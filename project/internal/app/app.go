package app

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"

	"github.com/labstack/echo/v4"
	"golang.org/x/sync/errgroup"

	pubSubRouter "tickets/internal/app/pubsub/router"
	"tickets/internal/provider/db"
	"tickets/internal/provider/outbox"
)

type Service struct {
	db           *db.DB
	httpRouter   *echo.Echo
	pubSubRouter *pubSubRouter.Router
	outbox       outbox.Outbox
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
