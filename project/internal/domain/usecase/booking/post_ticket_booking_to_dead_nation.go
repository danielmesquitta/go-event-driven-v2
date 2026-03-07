package booking

import (
	"context"
	"fmt"
	"tickets/internal/pkg/validator"
	"tickets/internal/provider/repo"
	"tickets/internal/provider/showapi"
)

type PostTicketBookingToDeadNation struct {
	showRepo repo.ShowRepo
	showAPI  showapi.ShowAPI
}

func NewPostTicketBookingToDeadNation(
	showRepo repo.ShowRepo,
	showAPI showapi.ShowAPI,
) *PostTicketBookingToDeadNation {
	return &PostTicketBookingToDeadNation{
		showRepo: showRepo,
		showAPI:  showAPI,
	}
}

type PostTicketBookingToDeadNationInput struct {
	BookingID       string `json:"booking_id" validate:"required"`
	ShowID          string `json:"show_id" validate:"required,uuid"`
	NumberOfTickets int    `json:"number_of_tickets" validate:"required"`
	CustomerEmail   string `json:"customer_email" validate:"required"`
}

func (uc *PostTicketBookingToDeadNation) Execute(
	ctx context.Context,
	in PostTicketBookingToDeadNationInput,
) error {
	err := validator.Validate(ctx, in)
	if err != nil {
		return fmt.Errorf("failed to validate input: %w", err)
	}

	show, err := uc.showRepo.Get(ctx, in.ShowID)
	if err != nil {
		return fmt.Errorf("failed to get show: %w", err)
	}
	if show == nil {
		return fmt.Errorf("show not found")
	}

	err = uc.showAPI.PostTicketBooking(ctx, showapi.Booking{
		BookingID:       in.BookingID,
		EventID:         show.DeadNationID,
		NumberOfTickets: in.NumberOfTickets,
		CustomerEmail:   in.CustomerEmail,
	})
	if err != nil {
		return fmt.Errorf("failed to post ticket booking: %w", err)
	}

	return nil
}
