package redisstream

import "github.com/ThreeDotsLabs/watermill/message"

func (p *PubSub) handleAppendToTracker(msg *message.Message) error {
	ticketID := string(msg.Payload)
	err := p.spreadsheetAPI.AppendRow(msg.Context(), "tickets-to-print", []string{ticketID})
	if err != nil {
		return err
	}
	return nil
}
