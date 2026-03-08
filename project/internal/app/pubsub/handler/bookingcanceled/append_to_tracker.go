package bookingcanceled

import (
	"context"
	"tickets/internal/app/pubsub/message/event"
	"tickets/internal/domain/usecase/tracker"
)

type AppendToTracker struct {
	appendToTrackerUseCase *tracker.AppendCanceledBooking
}

func NewAppendToTracker(
	appendToTrackerUseCase *tracker.AppendCanceledBooking,
) *AppendToTracker {
	return &AppendToTracker{
		appendToTrackerUseCase: appendToTrackerUseCase,
	}
}

func (a *AppendToTracker) Handle(ctx context.Context, event *event.TicketBookingCanceled) error {
	err := a.appendToTrackerUseCase.Execute(
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
