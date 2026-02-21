package handler

import (
	"context"
	"tickets/internal/app/pubsub/event"
	"tickets/internal/domain/usecase"
)

type AppendCanceledBookingToTracker struct {
	appendCanceledBookingToTrackerUseCase *usecase.AppendCanceledBookingToTracker
}

func NewAppendCanceledBookingToTracker(
	appendCanceledBookingToTrackerUseCase *usecase.AppendCanceledBookingToTracker,
) *AppendCanceledBookingToTracker {
	return &AppendCanceledBookingToTracker{
		appendCanceledBookingToTrackerUseCase: appendCanceledBookingToTrackerUseCase,
	}
}

func (a *AppendCanceledBookingToTracker) Handle(ctx context.Context, event *event.TicketBookingCanceled) error {
	err := a.appendCanceledBookingToTrackerUseCase.Execute(
		ctx,
		usecase.AppendCanceledBookingToTrackerInput{
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
