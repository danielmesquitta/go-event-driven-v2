package ticket

import (
	"context"
	"fmt"
	"tickets/internal/pkg/validator"
	"tickets/internal/provider/receipt"
)

type Refund struct {
	receiptsService receipt.Service
}

func NewRefund(
	receiptsService receipt.Service,
) *Refund {
	return &Refund{
		receiptsService: receiptsService,
	}
}

type RefundInput struct {
	TicketID string `json:"ticket_id" validate:"required"`
}

func (c *Refund) Execute(ctx context.Context, in RefundInput) error {
	err := validator.Validate(ctx, in)
	if err != nil {
		return fmt.Errorf("error validating ticket refund input: %w", err)
	}

	err = c.receiptsService.DeleteReceipt(ctx, in.TicketID, "customer requested refund")
	if err != nil {
		return fmt.Errorf("error deleting receipt: %w", err)
	}

	return nil
}
