package deadnation

import (
	"context"
	"fmt"
	"tickets/internal/provider/showapi"

	"github.com/ThreeDotsLabs/go-event-driven/v2/common/clients"
	"github.com/ThreeDotsLabs/go-event-driven/v2/common/clients/dead_nation"
	"github.com/google/uuid"
)

type DeadNationAPI struct {
	gateway *clients.Clients
}

func New(gateway *clients.Clients) *DeadNationAPI {
	return &DeadNationAPI{gateway: gateway}
}

func (d *DeadNationAPI) PostTicketBooking(ctx context.Context, booking showapi.Booking) error {
	eventID, err := uuid.Parse(booking.EventID)
	if err != nil {
		return fmt.Errorf("failed to parse show ID: %w", err)
	}

	bookingID, err := uuid.Parse(booking.BookingID)
	if err != nil {
		return fmt.Errorf("failed to parse booking ID: %w", err)
	}

	_, err = d.gateway.DeadNation.PostTicketBookingWithResponse(ctx, dead_nation.PostTicketBookingJSONRequestBody{
		BookingId:       bookingID,
		EventId:         eventID,
		NumberOfTickets: booking.NumberOfTickets,
		CustomerAddress: booking.CustomerEmail,
	}, nil)
	if err != nil {
		return fmt.Errorf("failed to get show: %w", err)
	}

	return nil
}

var _ showapi.ShowAPI = (*DeadNationAPI)(nil)
