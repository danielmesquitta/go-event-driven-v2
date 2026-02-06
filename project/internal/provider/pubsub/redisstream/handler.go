package redisstream

import (
	"context"
	"fmt"
	"tickets/internal/provider/pubsub"

	"github.com/ThreeDotsLabs/watermill/message"
)

func (p *PubSub) AddConsumerHandler(
	handlerName string,
	subscribeTopic string,
	handlerFunc message.NoPublishHandlerFunc,
) *message.Handler {
	sub, err := p.NewSubscriber(fmt.Sprintf("%s-%s-group", handlerName, subscribeTopic))
	if err != nil {
		panic(err)
	}

	return p.router.AddConsumerHandler(handlerName, subscribeTopic, sub, handlerFunc)
}

func (p *PubSub) RunRouter(
	ctx context.Context,
) error {
	p.AddConsumerHandler(
		"HandleIssueReceipt",
		pubsub.TopicIssueReceipt,
		p.handleIssueReceipt,
	)

	p.AddConsumerHandler(
		"HandleAppendToTracker",
		pubsub.TopicAppendToTracker,
		p.handleAppendToTracker,
	)

	return p.router.Run(ctx)
}

func (p *PubSub) handleIssueReceipt(msg *message.Message) error {
	ticketID := string(msg.Payload)
	err := p.receiptsService.IssueReceipt(msg.Context(), ticketID)
	if err != nil {
		return err
	}
	return nil
}

func (p *PubSub) handleAppendToTracker(msg *message.Message) error {
	ticketID := string(msg.Payload)
	err := p.spreadsheetAPI.AppendRow(msg.Context(), "tickets-to-print", []string{ticketID})
	if err != nil {
		return err
	}
	return nil
}
