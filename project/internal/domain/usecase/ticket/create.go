package ticket

import (
	"context"
	"tickets/internal/domain/entity"
	"tickets/internal/pkg/validator"
	"tickets/internal/provider/repo"
)

type Create struct {
	ticketRepo repo.TicketRepo
}

func NewCreate(ticketRepo repo.TicketRepo) *Create {
	return &Create{ticketRepo: ticketRepo}
}

type CreateInput struct {
	TicketID      string       `json:"ticket_id" validate:"required"`
	Price         entity.Money `json:"price" validate:"required"`
	CustomerEmail string       `json:"customer_email" validate:"required"`
}

func (c *Create) Execute(ctx context.Context, in CreateInput) error {
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
