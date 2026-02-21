package usecase

import (
	"context"
	"tickets/internal/domain/entity"
	"tickets/internal/provider/receipt"
)

type IssueReceipt struct {
	receiptsService receipt.Service
}

func NewIssueReceipt(receiptsService receipt.Service) *IssueReceipt {
	return &IssueReceipt{receiptsService: receiptsService}
}

type IssueReceiptInput struct {
	TicketID string
	Price    entity.Money
}

func (i *IssueReceipt) Execute(ctx context.Context, in IssueReceiptInput) error {
	err := i.receiptsService.IssueReceipt(
		ctx,
		in.TicketID,
		in.Price,
	)
	if err != nil {
		return err
	}
	return nil
}
