package redisstream

import "github.com/ThreeDotsLabs/watermill/message"

func (p *PubSub) handleIssueReceipt(msg *message.Message) error {
	ticketID := string(msg.Payload)
	err := p.receiptsService.IssueReceipt(msg.Context(), ticketID)
	if err != nil {
		return err
	}
	return nil
}
