package bookingmade

import (
	"context"
	"tickets/internal/app/pubsub/event"
	"tickets/internal/domain/usecase/booking"
)

type PostTicketBookingToDeadNation struct {
	mapShowIdToDeadNationEventIdUseCase *booking.PostTicketBookingToDeadNation
}

func NewPostTicketBookingToDeadNation(
	mapShowIdToDeadNationEventIdUseCase *booking.PostTicketBookingToDeadNation,
) *PostTicketBookingToDeadNation {
	return &PostTicketBookingToDeadNation{
		mapShowIdToDeadNationEventIdUseCase: mapShowIdToDeadNationEventIdUseCase,
	}
}

func (a *PostTicketBookingToDeadNation) Handle(ctx context.Context, event *event.BookingMade) error {
	err := a.mapShowIdToDeadNationEventIdUseCase.Execute(ctx, booking.PostTicketBookingToDeadNationInput{
		BookingID:       event.BookingID,
		ShowID:          event.ShowID,
		NumberOfTickets: event.NumberOfTickets,
		CustomerEmail:   event.CustomerEmail,
	})
	if err != nil {
		return err
	}
	return nil
}
