package usecase

import (
	"context"
	"fmt"
	"tickets/internal/app/pubsub/event"
	"tickets/internal/domain/entity"
	"tickets/internal/provider/eventbus"
	"tickets/internal/provider/filestorage"
)

type PrintTicket struct {
	fileStorage filestorage.Storage
	eventBus    eventbus.EventBus
}

func NewPrintTicket(
	fileStorage filestorage.Storage,
	eventBus eventbus.EventBus,
) *PrintTicket {
	return &PrintTicket{
		fileStorage: fileStorage,
		eventBus:    eventBus,
	}
}

type PrintTicketInput struct {
	CorrelationID string
	TicketID      string
	Price         entity.Money
}

func (p *PrintTicket) Execute(ctx context.Context, in PrintTicketInput) error {
	fileID := fmt.Sprintf("%s-ticket.html", in.TicketID)

	content := fmt.Sprintf(
		`<html><head><title>Ticket %s</title></head><body><h1>Ticket %s</h1><p>Price: %s %s</p></body></html>`,
		in.TicketID, in.TicketID, in.Price.Amount, in.Price.Currency,
	)

	err := p.fileStorage.StoreFile(ctx, fileID, content)
	if err != nil {
		return err
	}

	e := event.NewTicketPrinted(in.TicketID, fileID)
	if err := p.eventBus.Publish(ctx, e, eventbus.WithCorrelationID(in.CorrelationID)); err != nil {
		return err
	}

	return nil
}
