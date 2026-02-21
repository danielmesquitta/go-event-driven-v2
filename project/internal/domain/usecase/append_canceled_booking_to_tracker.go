package usecase

import (
	"context"
	"tickets/internal/domain/entity"
	"tickets/internal/provider/spreadsheet"
)

type AppendCanceledBookingToTracker struct {
	spreadsheetAPI spreadsheet.API
}

func NewAppendCanceledBookingToTracker(spreadsheetAPI spreadsheet.API) *AppendCanceledBookingToTracker {
	return &AppendCanceledBookingToTracker{spreadsheetAPI: spreadsheetAPI}
}

type AppendCanceledBookingToTrackerInput struct {
	TicketID      string
	CustomerEmail string
	Price         entity.Money
}

func (a *AppendCanceledBookingToTracker) Execute(ctx context.Context, in AppendCanceledBookingToTrackerInput) error {
	err := a.spreadsheetAPI.AppendRow(
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
