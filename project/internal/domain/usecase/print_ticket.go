package usecase

import (
	"context"
	"fmt"
	"tickets/internal/app/pubsub/event"
	"tickets/internal/domain/entity"
	"tickets/internal/pkg/ctxval"
	"tickets/internal/pkg/validator"
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
	TicketID string       `json:"ticket_id" validate:"required"`
	Price    entity.Money `json:"price" validate:"required"`
}

func (p *PrintTicket) Execute(ctx context.Context, in PrintTicketInput) error {
	err := validator.Validate(ctx, in)
	if err != nil {
		return err
	}

	fileID := fmt.Sprintf("%s-ticket.html", in.TicketID)

	content := fmt.Sprintf(
		`<html><head><title>Ticket %s</title></head><body><h1>Ticket %s</h1><p>Price: %s %s</p></body></html>`,
		in.TicketID, in.TicketID, in.Price.Amount, in.Price.Currency,
	)

	err = p.fileStorage.StoreFile(ctx, fileID, content)
	if err != nil {
		return err
	}

	idempotencyKey := ctxval.GetIdempotencyKey(ctx) + "-" + in.TicketID
	ctx = ctxval.WithIdempotencyKey(ctx, idempotencyKey)

	e := event.NewTicketPrinted(
		ctx,
		in.TicketID,
		fileID,
	)
	if err := p.eventBus.Publish(ctx, e); err != nil {
		return err
	}

	return nil
}
