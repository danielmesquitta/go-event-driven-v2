package handler

import (
	"context"
	"tickets/internal/app/pubsub/event"
	"tickets/internal/domain/usecase/tracker"
)

type AppendCanceledBookingToTracker struct {
	appendCanceledBookingToTrackerUseCase *tracker.AppendCanceledBooking
}

func NewAppendCanceledBookingToTracker(
	appendCanceledBookingToTrackerUseCase *tracker.AppendCanceledBooking,
) *AppendCanceledBookingToTracker {
	return &AppendCanceledBookingToTracker{
		appendCanceledBookingToTrackerUseCase: appendCanceledBookingToTrackerUseCase,
	}
}

func (a *AppendCanceledBookingToTracker) Handle(ctx context.Context, event *event.TicketBookingCanceled) error {
	err := a.appendCanceledBookingToTrackerUseCase.Execute(
		ctx,
		tracker.AppendCanceledBookingInput{
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
