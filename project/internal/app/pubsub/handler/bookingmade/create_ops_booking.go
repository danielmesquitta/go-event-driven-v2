package bookingmade

import (
	"context"
	"tickets/internal/app/pubsub/message/event"
	"tickets/internal/domain/usecase/opsbooking"
)

type CreateOpsBooking struct {
	createOpsBookingUseCase *opsbooking.Create
}

func NewCreateOpsBooking(
	createOpsBookingUseCase *opsbooking.Create,
) *CreateOpsBooking {
	return &CreateOpsBooking{
		createOpsBookingUseCase: createOpsBookingUseCase,
	}
}

func (a *CreateOpsBooking) Handle(ctx context.Context, event *event.BookingMade) error {
	err := a.createOpsBookingUseCase.Execute(ctx, opsbooking.CreateInput{
		BookingID: event.BookingID,
		BookedAt:  event.Header.PublishedAt.UTC(),
	})
	if err != nil {
		return err
	}
	return nil
}
