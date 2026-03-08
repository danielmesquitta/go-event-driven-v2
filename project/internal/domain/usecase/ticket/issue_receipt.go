package ticket

import (
	"context"
	"fmt"
	"tickets/internal/app/pubsub/message/event"
	"tickets/internal/domain/entity"
	"tickets/internal/pkg/validator"
	"tickets/internal/provider/eventbus"
	"tickets/internal/provider/receipt"
)

type IssueReceipt struct {
	receiptsService receipt.Service
	eventBus        eventbus.EventBus
}

func NewIssueReceipt(receiptsService receipt.Service, eventBus eventbus.EventBus) *IssueReceipt {
	return &IssueReceipt{receiptsService: receiptsService, eventBus: eventBus}
}

type IssueReceiptInput struct {
	TicketID string       `json:"ticket_id" validate:"required"`
	Price    entity.Money `json:"price" validate:"required"`
}

func (i *IssueReceipt) Execute(ctx context.Context, in IssueReceiptInput) error {
	err := validator.Validate(ctx, in)
	if err != nil {
		return fmt.Errorf("failed to validate input: %w", err)
	}

	receipt, err := i.receiptsService.IssueReceipt(
		ctx,
		in.TicketID,
		in.Price,
	)
	if err != nil {
		return fmt.Errorf("failed to issue receipt: %w", err)
	}

	e := event.NewTicketReceiptIssued(
		ctx,
		in.TicketID,
		receipt.ReceiptNumber,
		receipt.IssuedAt,
	)
	if err := i.eventBus.Publish(ctx, e); err != nil {
		return fmt.Errorf("failed to publish TicketReceiptIssued event: %w", err)
	}

	return nil
}
