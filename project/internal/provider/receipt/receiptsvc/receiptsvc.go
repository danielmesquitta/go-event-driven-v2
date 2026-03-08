package receiptsvc

import (
	"context"
	"encoding/json"
	"fmt"
	"tickets/internal/domain/entity"
	"tickets/internal/pkg/ctxval"
	"tickets/internal/provider/receipt"
	"time"

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

type IssueReceiptResponse struct {
	Number   string    `json:"number"`
	IssuedAt time.Time `json:"issued_at"`
}

func (c Client) IssueReceipt(ctx context.Context, ticketID string, price entity.Money) (*receipt.Receipt, error) {
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
		return nil, fmt.Errorf("failed to post receipt: %w", err)
	}
	if resp.StatusCode() < 200 || resp.StatusCode() >= 300 {
		return nil, fmt.Errorf("unexpected status code for IssueReceipt: %d", resp.StatusCode())
	}

	var issueReceiptResponse IssueReceiptResponse
	err = json.Unmarshal(resp.Body, &issueReceiptResponse)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal IssueReceiptResponse: %w", err)
	}

	receipt := receipt.Receipt{
		ReceiptNumber: issueReceiptResponse.Number,
		IssuedAt:      issueReceiptResponse.IssuedAt,
	}

	return &receipt, nil
}

func (c Client) DeleteReceipt(ctx context.Context, ticketID string, reason string) error {
	idempotencyKey := ctxval.GetIdempotencyKey(ctx)
	resp, err := c.clients.Receipts.PutVoidReceiptWithResponse(ctx, receipts.VoidReceiptRequest{
		IdempotentId: &idempotencyKey,
		TicketId:     ticketID,
		Reason:       reason,
	})
	if err != nil {
		return fmt.Errorf("failed to post receipt: %w", err)
	}

	if resp.StatusCode() < 200 || resp.StatusCode() >= 300 {
		return fmt.Errorf("unexpected status code for DeleteReceipt: %d", resp.StatusCode())
	}

	return nil
}

var _ receipt.Service = &Client{}
