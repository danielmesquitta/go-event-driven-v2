package pubsub

import (
	"encoding/json"
	"tickets/internal/app/pubsub/event"

	"github.com/ThreeDotsLabs/watermill/message"
)

func (r *Router) handleIssueReceipt(msg *message.Message) error {
	var issueReceiptEvent event.IssueReceiptEvent
	err := json.Unmarshal(msg.Payload, &issueReceiptEvent)
	if err != nil {
		return err
	}

	err = r.receiptsService.IssueReceipt(
		msg.Context(),
		issueReceiptEvent.TicketID,
		issueReceiptEvent.Price,
	)
	if err != nil {
		return err
	}
	return nil
}
