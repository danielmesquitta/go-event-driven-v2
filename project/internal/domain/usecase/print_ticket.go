package usecase

import (
	"context"
	"fmt"
	"tickets/internal/domain/entity"
	"tickets/internal/provider/filestorage"
)

type PrintTicket struct {
	fileStorage filestorage.Storage
}

func NewPrintTicket(fileStorage filestorage.Storage) *PrintTicket {
	return &PrintTicket{fileStorage: fileStorage}
}

type PrintTicketInput struct {
	TicketID string
	Price    entity.Money
}

func (p *PrintTicket) Execute(ctx context.Context, in PrintTicketInput) error {
	fileID := fmt.Sprintf("%s-ticket.html", in.TicketID)

	content := fmt.Sprintf(
		`<html><head><title>Ticket %s</title></head><body><h1>Ticket %s</h1><p>Price: %s %s</p></body></html>`,
		in.TicketID, in.TicketID, in.Price.Amount, in.Price.Currency,
	)

	return p.fileStorage.StoreFile(ctx, fileID, content)
}
