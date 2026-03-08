package receipt

import (
	"context"
	"tickets/internal/domain/entity"
	"time"
)

type Receipt struct {
	ReceiptNumber string
	IssuedAt      time.Time
}

type Service interface {
	IssueReceipt(ctx context.Context, ticketID string, price entity.Money) (*Receipt, error)
	DeleteReceipt(ctx context.Context, ticketID string, reason string) error
}
