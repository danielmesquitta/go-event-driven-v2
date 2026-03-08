package paymentsvc

import (
	"context"
	"fmt"
	"tickets/internal/pkg/ctxval"
	"tickets/internal/provider/payment"

	"github.com/ThreeDotsLabs/go-event-driven/v2/common/clients"
	"github.com/ThreeDotsLabs/go-event-driven/v2/common/clients/payments"
)

type Client struct {
	clients *clients.Clients
}

func NewClient(clients *clients.Clients) *Client {
	if clients == nil {
		panic("NewClient: clients is nil")
	}

	return &Client{clients: clients}
}

func (c Client) Refund(ctx context.Context, ticketID string, reason string) error {
	idempotencyKey := ctxval.GetIdempotencyKey(ctx)
	resp, err := c.clients.Payments.PutRefundsWithResponse(ctx, payments.PaymentRefundRequest{
		DeduplicationId:  &idempotencyKey,
		Reason:           reason,
		PaymentReference: ticketID,
	})
	if err != nil {
		return fmt.Errorf("failed to post refund: %w", err)
	}

	if resp.StatusCode() < 200 || resp.StatusCode() >= 300 {
		return fmt.Errorf("unexpected status code for PutRefundsWithResponse: %d", resp.StatusCode())
	}

	return nil
}

var _ payment.Service = &Client{}
