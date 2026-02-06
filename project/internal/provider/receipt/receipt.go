package receipt

import "context"

type Service interface {
	IssueReceipt(ctx context.Context, ticketID string) error
}
