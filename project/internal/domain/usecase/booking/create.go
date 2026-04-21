package booking

import (
	"context"
	"fmt"
	"tickets/internal/app/pubsub/message/event"
	"tickets/internal/domain/entity"
	"tickets/internal/domain/errs"
	"tickets/internal/pkg/tx"
	"tickets/internal/pkg/validator"
	"tickets/internal/provider/outbox"
	"tickets/internal/provider/repo"

	"github.com/google/uuid"
)

type Create struct {
	tx          tx.Transaction
	showRepo    repo.ShowRepo
	bookingRepo repo.BookingRepo
	outbox      outbox.Outbox
}

func NewCreate(
	tx tx.Transaction,
	showRepo repo.ShowRepo,
	bookingRepo repo.BookingRepo,
	outbox outbox.Outbox,
) *Create {
	return &Create{
		tx:          tx,
		showRepo:    showRepo,
		bookingRepo: bookingRepo,
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
		return "", fmt.Errorf("failed to validate input: %w", err)
	}

	availableTickets, err := uc.showRepo.AvailableTickets(ctx, in.ShowID)
	if err != nil {
		return "", fmt.Errorf("failed to get available tickets: %w", err)
	}
	if availableTickets < in.NumberOfTickets {
		return "", errs.ErrBadRequest.New(
			errs.WithMessage("not enough available tickets"),
			errs.WithMetadata(errs.MetadataDataKey, map[string]int{
				"available_tickets": availableTickets,
				"requested_tickets": in.NumberOfTickets,
			}),
		)
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
