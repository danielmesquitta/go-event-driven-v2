package usecase

import (
	"context"
	"tickets/internal/domain/entity"
	"tickets/internal/pkg/validator"
	"tickets/internal/provider/spreadsheet"
)

type AppendCanceledBookingToTracker struct {
	spreadsheetAPI spreadsheet.API
}

func NewAppendCanceledBookingToTracker(spreadsheetAPI spreadsheet.API) *AppendCanceledBookingToTracker {
	return &AppendCanceledBookingToTracker{spreadsheetAPI: spreadsheetAPI}
}

type AppendCanceledBookingToTrackerInput struct {
	TicketID      string       `json:"ticket_id" validate:"required"`
	CustomerEmail string       `json:"customer_email" validate:"required"`
	Price         entity.Money `json:"price" validate:"required"`
}

func (a *AppendCanceledBookingToTracker) Execute(ctx context.Context, in AppendCanceledBookingToTrackerInput) error {
	err := validator.Validate(ctx, in)
	if err != nil {
		return err
	}

	err = a.spreadsheetAPI.AppendRow(
		ctx,
		"tickets-to-refund",
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
