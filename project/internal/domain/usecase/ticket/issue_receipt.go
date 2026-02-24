package ticket

import (
	"context"
	"tickets/internal/domain/entity"
	"tickets/internal/pkg/validator"
	"tickets/internal/provider/receipt"
)

type IssueReceipt struct {
	receiptsService receipt.Service
}

func NewIssueReceipt(receiptsService receipt.Service) *IssueReceipt {
	return &IssueReceipt{receiptsService: receiptsService}
}

type IssueReceiptInput struct {
	TicketID string       `json:"ticket_id" validate:"required"`
	Price    entity.Money `json:"price" validate:"required"`
}

func (i *IssueReceipt) Execute(ctx context.Context, in IssueReceiptInput) error {
	err := validator.Validate(ctx, in)
	if err != nil {
		return err
	}

	err = i.receiptsService.IssueReceipt(
		ctx,
		in.TicketID,
		in.Price,
	)
	if err != nil {
		return err
	}
	return nil
}
