package pubsub

import "github.com/ThreeDotsLabs/watermill/message"

func (r *Router) handleIssueReceipt(msg *message.Message) error {
	ticketID := string(msg.Payload)
	err := r.receiptsService.IssueReceipt(msg.Context(), ticketID)
	if err != nil {
		return err
	}
	return nil
}
