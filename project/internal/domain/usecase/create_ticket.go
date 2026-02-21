package usecase

import (
	"context"
	"tickets/internal/domain/entity"
	"tickets/internal/provider/repo"
)

type CreateTicket struct {
	ticketRepo repo.TicketRepo
}

func NewCreateTicket(ticketRepo repo.TicketRepo) *CreateTicket {
	return &CreateTicket{ticketRepo: ticketRepo}
}

type CreateTicketInput struct {
	TicketID      string
	Price         entity.Money
	CustomerEmail string
}

func (c *CreateTicket) Execute(ctx context.Context, in CreateTicketInput) error {
	err := c.ticketRepo.Create(ctx, &entity.Ticket{
		ID:            in.TicketID,
		Price:         in.Price,
		CustomerEmail: in.CustomerEmail,
	})
	if err != nil {
		return err
	}
	return nil
}
