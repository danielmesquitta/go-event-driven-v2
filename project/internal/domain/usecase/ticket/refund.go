package ticket

import (
	"context"
	"fmt"
	"tickets/internal/pkg/validator"
	"tickets/internal/provider/payment"
	"tickets/internal/provider/receipt"

	"golang.org/x/sync/errgroup"
)

const reason = "customer requested refund"

type Refund struct {
	receiptsService receipt.Service
	paymentsService payment.Service
}

func NewRefund(
	receiptsService receipt.Service,
	paymentsService payment.Service,
) *Refund {
	return &Refund{
		receiptsService: receiptsService,
		paymentsService: paymentsService,
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

	g, gCtx := errgroup.WithContext(ctx)

	g.Go(func() error {
		return c.receiptsService.DeleteReceipt(gCtx, in.TicketID, reason)
	})

	g.Go(func() error {
		return c.paymentsService.Refund(gCtx, in.TicketID, reason)
	})

	if err := g.Wait(); err != nil {
		return fmt.Errorf("error refunding ticket: %w", err)
	}

	return nil
}
