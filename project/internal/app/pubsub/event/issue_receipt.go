package event

import "tickets/internal/domain/entity"

const TopicIssueReceipt Topic = "issue-receipt"

type IssueReceiptEvent struct {
	TicketID string       `json:"ticket_id"`
	Price    entity.Money `json:"price"`
}
