package pubsub

import (
	"encoding/json"
	"tickets/internal/app/pubsub/event"

	"github.com/ThreeDotsLabs/watermill/message"
)

func (r *Router) handleAppendToTracker(msg *message.Message) error {
	var event event.AppendToTrackerEvent
	err := json.Unmarshal(msg.Payload, &event)
	if err != nil {
		return err
	}

	err = r.spreadsheetAPI.AppendRow(
		msg.Context(),
		"tickets-to-print",
		[]string{event.TicketID, event.CustomerEmail, event.Price.Amount, event.Price.Currency},
	)
	if err != nil {
		return err
	}
	return nil
}
