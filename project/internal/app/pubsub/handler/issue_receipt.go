package handler

import (
	"encoding/json"
	"tickets/internal/app/pubsub/event"
	"tickets/internal/domain/errs"
	"tickets/internal/provider/receipt"

	"github.com/ThreeDotsLabs/watermill/message"
)

type IssueReceipt struct {
	receiptsService receipt.Service
}

func NewIssueReceipt(receiptsService receipt.Service) *IssueReceipt {
	return &IssueReceipt{receiptsService: receiptsService}
}

func (i *IssueReceipt) Handle(msg *message.Message) error {
	var e event.TicketBookingConfirmed
	err := json.Unmarshal(msg.Payload, &e)
	if err != nil {
		return errs.ErrInvalidFormat.New(errs.WithMetadata("payload", string(msg.Payload)))
	}

	e.SetDefaults()

	err = i.receiptsService.IssueReceipt(
		msg.Context(),
		e.TicketID,
		e.Price,
	)
	if err != nil {
		return err
	}
	return nil
}
