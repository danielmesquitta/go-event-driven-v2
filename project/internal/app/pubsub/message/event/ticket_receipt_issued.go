package event

import (
	"context"
	"tickets/internal/app/pubsub/message"
	"time"
)

type TicketReceiptIssued struct {
	Header        message.Header `json:"header"`
	TicketID      string         `json:"ticket_id"`
	ReceiptNumber string         `json:"receipt_number"`
	IssuedAt      time.Time      `json:"issued_at"`
}

func NewTicketReceiptIssued(
	ctx context.Context,
	ticketID string,
	receiptNumber string,
	issuedAt time.Time,
) *TicketReceiptIssued {
	return &TicketReceiptIssued{
		Header:        message.NewHeader(ctx),
		TicketID:      ticketID,
		ReceiptNumber: receiptNumber,
		IssuedAt:      issuedAt,
	}
}

func (e *TicketReceiptIssued) GetHeader() message.Header {
	return e.Header
}

var _ Event = (*TicketReceiptIssued)(nil)
