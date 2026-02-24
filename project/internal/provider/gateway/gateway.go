package gateway

import (
	"context"
	"net/http"
	"os"
	"tickets/internal/pkg/ctxval"

	"github.com/ThreeDotsLabs/go-event-driven/v2/common/clients"
)

func NewGateway() *clients.Clients {
	clients, err := clients.NewClients(
		os.Getenv("GATEWAY_ADDR"),
		func(ctx context.Context, req *http.Request) error {
			req.Header.Set("Correlation-ID", ctxval.GetCorrelationID(ctx))
			return nil
		},
	)
	if err != nil {
		panic(err)
	}

	return clients
}
