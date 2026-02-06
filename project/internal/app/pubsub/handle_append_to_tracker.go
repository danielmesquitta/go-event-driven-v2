package pubsub

import "github.com/ThreeDotsLabs/watermill/message"

func (r *Router) handleAppendToTracker(msg *message.Message) error {
	ticketID := string(msg.Payload)
	err := r.spreadsheetAPI.AppendRow(msg.Context(), "tickets-to-print", []string{ticketID})
	if err != nil {
		return err
	}
	return nil
}
