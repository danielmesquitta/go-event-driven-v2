package bookingconfirmed

import (
	"context"
	"tickets/internal/app/pubsub/message/event"
	"tickets/internal/domain/usecase/tracker"
)

type AppendToTracker struct {
	appendToTrackerUseCase *tracker.AppendConfirmedBooking
}

func NewAppendToTracker(
	appendToTrackerUseCase *tracker.AppendConfirmedBooking,
) *AppendToTracker {
	return &AppendToTracker{
		appendToTrackerUseCase: appendToTrackerUseCase,
	}
}

func (a *AppendToTracker) Handle(ctx context.Context, event *event.TicketBookingConfirmed) error {
	err := a.appendToTrackerUseCase.Execute(
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
