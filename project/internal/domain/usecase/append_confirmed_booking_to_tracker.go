package usecase

import (
	"context"
	"tickets/internal/domain/entity"
	"tickets/internal/pkg/validator"
	"tickets/internal/provider/spreadsheet"
)

type AppendConfirmedBookingToTracker struct {
	spreadsheetAPI spreadsheet.API
}

func NewAppendConfirmedBookingToTracker(spreadsheetAPI spreadsheet.API) *AppendConfirmedBookingToTracker {
	return &AppendConfirmedBookingToTracker{spreadsheetAPI: spreadsheetAPI}
}

type AppendConfirmedBookingToTrackerInput struct {
	TicketID      string       `json:"ticket_id" validate:"required"`
	CustomerEmail string       `json:"customer_email" validate:"required"`
	Price         entity.Money `json:"price" validate:"required"`
}

func (a *AppendConfirmedBookingToTracker) Execute(ctx context.Context, in AppendConfirmedBookingToTrackerInput) error {
	err := validator.Validate(ctx, in)
	if err != nil {
		return err
	}

	err = a.spreadsheetAPI.AppendRow(
		ctx,
		"tickets-to-print",
		[]string{
			in.TicketID,
			in.CustomerEmail,
			in.Price.Amount,
			in.Price.Currency,
		},
	)
	if err != nil {
		return err
	}
	return nil
}
