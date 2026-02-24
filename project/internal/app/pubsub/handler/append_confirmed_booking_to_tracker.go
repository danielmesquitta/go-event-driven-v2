package handler

import (
	"context"
	"tickets/internal/app/pubsub/event"
	"tickets/internal/domain/usecase/tracker"
)

type AppendConfirmedBookingToTracker struct {
	appendConfirmedBookingToTrackerUseCase *tracker.AppendConfirmedBooking
}

func NewAppendConfirmedBookingToTracker(
	appendConfirmedBookingToTrackerUseCase *tracker.AppendConfirmedBooking,
) *AppendConfirmedBookingToTracker {
	return &AppendConfirmedBookingToTracker{
		appendConfirmedBookingToTrackerUseCase: appendConfirmedBookingToTrackerUseCase,
	}
}

func (a *AppendConfirmedBookingToTracker) Handle(ctx context.Context, event *event.TicketBookingConfirmed) error {
	err := a.appendConfirmedBookingToTrackerUseCase.Execute(
		ctx,
		tracker.AppendConfirmedBookingInput{
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
