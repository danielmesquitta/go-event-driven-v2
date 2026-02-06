package pubsub

import (
	"encoding/json"
	"tickets/internal/app/pubsub/event"

	"github.com/ThreeDotsLabs/watermill/message"
)

func (r *Router) handleAppendToTracker(msg *message.Message) error {
	var e event.TicketBookingConfirmed
	err := json.Unmarshal(msg.Payload, &e)
	if err != nil {
		return err
	}

	err = r.spreadsheetAPI.AppendRow(
		msg.Context(),
		"tickets-to-print",
		[]string{
			e.TicketID,
			e.CustomerEmail,
			e.Price.Amount,
			e.Price.Currency,
		},
	)
	if err != nil {
		return err
	}
	return nil
}
