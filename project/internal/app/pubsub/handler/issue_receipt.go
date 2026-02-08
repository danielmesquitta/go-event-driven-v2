package handler

import (
	"context"
	"tickets/internal/app/pubsub/event"
	"tickets/internal/provider/receipt"
)

type IssueReceipt struct {
	receiptsService receipt.Service
}

func NewIssueReceipt(receiptsService receipt.Service) *IssueReceipt {
	return &IssueReceipt{receiptsService: receiptsService}
}

func (i *IssueReceipt) Handle(ctx context.Context, event *event.TicketBookingConfirmed) error {
	err := i.receiptsService.IssueReceipt(
		ctx,
		event.TicketID,
		event.Price,
	)
	if err != nil {
		return err
	}
	return nil
}
