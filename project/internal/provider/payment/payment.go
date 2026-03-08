package payment

import (
	"context"
)

type Service interface {
	Refund(ctx context.Context, ticketID string, reason string) error
}
