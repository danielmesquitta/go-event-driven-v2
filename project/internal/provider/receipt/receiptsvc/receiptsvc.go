package receiptsvc

import (
	"context"
	"fmt"
	"net/http"
	"tickets/internal/domain/entity"
	"tickets/internal/pkg/ctxval"
	"tickets/internal/provider/receipt"

	"github.com/ThreeDotsLabs/go-event-driven/v2/common/clients"
	"github.com/ThreeDotsLabs/go-event-driven/v2/common/clients/receipts"
)

type Client struct {
	// we are not mocking this client: it's pointless to use interface here
	clients *clients.Clients
}

func NewClient(clients *clients.Clients) *Client {
	if clients == nil {
		panic("NewClient: clients is nil")
	}

	return &Client{clients: clients}
}

func (c Client) IssueReceipt(ctx context.Context, ticketID string, price entity.Money) error {
	idempotencyKey := ctxval.GetIdempotencyKey(ctx)
	resp, err := c.clients.Receipts.PutReceiptsWithResponse(ctx, receipts.CreateReceipt{
		IdempotencyKey: &idempotencyKey,
		TicketId:       ticketID,
		Price: receipts.Money{
			MoneyAmount:   price.Amount,
			MoneyCurrency: price.Currency,
		},
	})
	if err != nil {
		return fmt.Errorf("failed to post receipt: %w", err)
	}

	switch resp.StatusCode() {
	case http.StatusOK:
		// receipt already exists
		return nil
	case http.StatusCreated:
		// receipt was created
		return nil
	default:
		return fmt.Errorf("unexpected status code for POST receipts-api/receipts: %d", resp.StatusCode())
	}
}

var _ receipt.Service = &Client{}
