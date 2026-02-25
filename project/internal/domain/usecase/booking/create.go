package booking

import (
	"context"
	"tickets/internal/domain/entity"
	"tickets/internal/pkg/validator"
	"tickets/internal/provider/repo"

	"github.com/google/uuid"
)

type Create struct {
	bookingRepo repo.BookingRepo
}

func NewCreate(bookingRepo repo.BookingRepo) *Create {
	return &Create{bookingRepo: bookingRepo}
}

type CreateInput struct {
	ShowID          string `json:"show_id" validate:"required"`
	NumberOfTickets int    `json:"number_of_tickets" validate:"required"`
	CustomerEmail   string `json:"customer_email" validate:"required"`
}

func (uc *Create) Execute(ctx context.Context, in CreateInput) (id string, err error) {
	err = validator.Validate(ctx, in)
	if err != nil {
		return "", err
	}

	id = uuid.NewString()

	err = uc.bookingRepo.Create(ctx, &entity.Booking{
		ID:              id,
		ShowID:          in.ShowID,
		NumberOfTickets: in.NumberOfTickets,
		CustomerEmail:   in.CustomerEmail,
	})
	if err != nil {
		return "", err
	}

	return id, nil
}
