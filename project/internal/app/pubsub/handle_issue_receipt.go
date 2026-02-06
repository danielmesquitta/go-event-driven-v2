package pubsub

import (
	"encoding/json"
	"tickets/internal/app/pubsub/event"

	"github.com/ThreeDotsLabs/watermill/message"
)

func (r *Router) handleIssueReceipt(msg *message.Message) error {
	var e event.TicketBookingConfirmed
	err := json.Unmarshal(msg.Payload, &e)
	if err != nil {
		return err
	}

	err = r.receiptsService.IssueReceipt(
		msg.Context(),
		e.TicketID,
		e.Price,
	)
	if err != nil {
		return err
	}
	return nil
}
