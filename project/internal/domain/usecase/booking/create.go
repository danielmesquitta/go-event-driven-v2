package booking

import (
	"context"
	"fmt"
	"tickets/internal/app/pubsub/event"
	"tickets/internal/domain/entity"
	"tickets/internal/pkg/tx"
	"tickets/internal/pkg/validator"
	"tickets/internal/provider/db"
	"tickets/internal/provider/outbox"
	"tickets/internal/provider/repo"

	"github.com/google/uuid"
)

type Create struct {
	tx          tx.Transaction
	bookingRepo repo.BookingRepo
	db          *db.DB
	outbox      outbox.Outbox
}

func NewCreate(
	tx tx.Transaction,
	bookingRepo repo.BookingRepo,
	database *db.DB,
	outbox outbox.Outbox,
) *Create {
	return &Create{
		tx:          tx,
		bookingRepo: bookingRepo,
		db:          database,
		outbox:      outbox,
	}
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

	err = uc.tx.Do(ctx, func(ctx context.Context) error {
		err = uc.bookingRepo.Create(ctx, &entity.Booking{
			ID:              id,
			ShowID:          in.ShowID,
			NumberOfTickets: in.NumberOfTickets,
			CustomerEmail:   in.CustomerEmail,
		})
		if err != nil {
			return fmt.Errorf("failed to create booking: %w", err)
		}

		evnt := event.NewBookingMade(ctx, id, in.NumberOfTickets, in.CustomerEmail, in.ShowID)
		if err := uc.outbox.Publish(ctx, evnt); err != nil {
			return fmt.Errorf("failed to publish BookingMade event: %w", err)
		}

		return nil
	})
	if err != nil {
		return "", fmt.Errorf("failed to transaction while creating booking: %w", err)
	}

	return id, nil
}
