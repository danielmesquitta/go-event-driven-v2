package handler

import (
	"context"
	"tickets/internal/app/pubsub/event"
	"tickets/internal/domain/usecase"
)

type IssueReceipt struct {
	issueReceiptUseCase *usecase.IssueReceipt
}

func NewIssueReceipt(issueReceiptUseCase *usecase.IssueReceipt) *IssueReceipt {
	return &IssueReceipt{issueReceiptUseCase: issueReceiptUseCase}
}

func (i *IssueReceipt) Handle(ctx context.Context, event *event.TicketBookingConfirmed) error {
	err := i.issueReceiptUseCase.Execute(ctx, usecase.IssueReceiptInput{
		TicketID: event.TicketID,
		Price:    event.Price,
	})
	if err != nil {
		return err
	}
	return nil
}
