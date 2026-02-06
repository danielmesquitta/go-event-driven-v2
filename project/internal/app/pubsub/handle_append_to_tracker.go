package pubsub

import (
	"encoding/json"
	"tickets/internal/app/pubsub/event"

	"github.com/ThreeDotsLabs/watermill/message"
)

func (r *Router) handleAppendToTracker(msg *message.Message) error {
	var appendToTrackerEvent event.AppendToTrackerEvent
	err := json.Unmarshal(msg.Payload, &appendToTrackerEvent)
	if err != nil {
		return err
	}

	err = r.spreadsheetAPI.AppendRow(
		msg.Context(),
		"tickets-to-print",
		[]string{
			appendToTrackerEvent.TicketID,
			appendToTrackerEvent.CustomerEmail,
			appendToTrackerEvent.Price.Amount,
			appendToTrackerEvent.Price.Currency,
		},
	)
	if err != nil {
		return err
	}
	return nil
}
