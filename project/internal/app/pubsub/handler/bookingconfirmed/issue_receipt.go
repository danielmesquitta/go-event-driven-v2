package bookingconfirmed

import (
	"context"
	"tickets/internal/app/pubsub/message/event"
	"tickets/internal/domain/usecase/ticket"
)

type IssueReceipt struct {
	issueReceiptUseCase *ticket.IssueReceipt
}

func NewIssueReceipt(issueReceiptUseCase *ticket.IssueReceipt) *IssueReceipt {
	return &IssueReceipt{issueReceiptUseCase: issueReceiptUseCase}
}

func (i *IssueReceipt) Handle(ctx context.Context, event *event.TicketBookingConfirmed) error {
	err := i.issueReceiptUseCase.Execute(ctx, ticket.IssueReceiptInput{
		TicketID: event.TicketID,
		Price:    event.Price,
	})
	if err != nil {
		return err
	}
	return nil
}
