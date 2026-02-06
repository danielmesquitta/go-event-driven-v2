package receipt

import (
	"context"
	"tickets/internal/domain/entity"
)

type Service interface {
	IssueReceipt(ctx context.Context, ticketID string, price entity.Money) error
}
