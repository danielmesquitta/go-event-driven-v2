package redisstream

import (
	"context"
	"log/slog"
	"tickets/internal/provider/pubsub"
)

func (p *PubSub) RunHandlers(
	ctx context.Context,
) error {
	if err := p.handleIssueReceipt(); err != nil {
		return err
	}

	if err := p.handleAppendToTracker(); err != nil {
		return err
	}

	return nil
}

func (p *PubSub) handleIssueReceipt() error {
	issueReceiptSub, err := p.NewSubscriber("issue-receipt-group")
	if err != nil {
		return err
	}

	go func() {
		messages, err := issueReceiptSub.Subscribe(context.Background(), pubsub.TopicIssueReceipt)
		if err != nil {
			panic(err)
		}

		for msg := range messages {
			ticketID := string(msg.Payload)
			if err := p.receiptsService.IssueReceipt(msg.Context(), ticketID); err != nil {
				slog.With("error", err).Error("failed to issue the receipt")
				msg.Nack()
			} else {
				msg.Ack()
			}
		}
	}()

	return nil
}

func (p *PubSub) handleAppendToTracker() error {
	appendToTrackerSub, err := p.NewSubscriber("append-to-tracker-group")
	if err != nil {
		return err
	}

	go func() {
		messages, err := appendToTrackerSub.Subscribe(context.Background(), pubsub.TopicAppendToTracker)
		if err != nil {
			panic(err)
		}

		for msg := range messages {
			ticketID := string(msg.Payload)
			if err := p.spreadsheetAPI.AppendRow(msg.Context(), "tickets-to-print", []string{ticketID}); err != nil {
				slog.With("error", err).Error("failed to append to tracker")
				msg.Nack()
			} else {
				msg.Ack()
			}
		}
	}()

	return nil
}
