package handler

import (
	"context"
	"tickets/internal/app/pubsub/event"
	"tickets/internal/domain/usecase"
)

type AppendConfirmedBookingToTracker struct {
	appendConfirmedBookingToTrackerUseCase *usecase.AppendConfirmedBookingToTracker
}

func NewAppendConfirmedBookingToTracker(
	appendConfirmedBookingToTrackerUseCase *usecase.AppendConfirmedBookingToTracker,
) *AppendConfirmedBookingToTracker {
	return &AppendConfirmedBookingToTracker{
		appendConfirmedBookingToTrackerUseCase: appendConfirmedBookingToTrackerUseCase,
	}
}

func (a *AppendConfirmedBookingToTracker) Handle(ctx context.Context, event *event.TicketBookingConfirmed) error {
	err := a.appendConfirmedBookingToTrackerUseCase.Execute(
		ctx,
		usecase.AppendConfirmedBookingToTrackerInput{
			TicketID:      event.TicketID,
			CustomerEmail: event.CustomerEmail,
			Price:         event.Price,
		},
	)
	if err != nil {
		return err
	}
	return nil
}
