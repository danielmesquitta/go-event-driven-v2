package usecase

import (
	"context"
	"tickets/internal/domain/entity"
	"tickets/internal/pkg/validator"
	"tickets/internal/provider/repo"
)

type CreateTicket struct {
	ticketRepo repo.TicketRepo
}

func NewCreateTicket(ticketRepo repo.TicketRepo) *CreateTicket {
	return &CreateTicket{ticketRepo: ticketRepo}
}

type CreateTicketInput struct {
	TicketID      string       `json:"ticket_id" validate:"required"`
	Price         entity.Money `json:"price" validate:"required"`
	CustomerEmail string       `json:"customer_email" validate:"required"`
}

func (c *CreateTicket) Execute(ctx context.Context, in CreateTicketInput) error {
	err := validator.Validate(ctx, in)
	if err != nil {
		return err
	}

	err = c.ticketRepo.Create(ctx, &entity.Ticket{
		ID:            in.TicketID,
		Price:         in.Price,
		CustomerEmail: in.CustomerEmail,
	})
	if err != nil {
		return err
	}
	return nil
}
