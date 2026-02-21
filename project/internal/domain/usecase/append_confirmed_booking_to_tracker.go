package usecase

import (
	"context"
	"tickets/internal/domain/entity"
	"tickets/internal/provider/spreadsheet"
)

type AppendConfirmedBookingToTracker struct {
	spreadsheetAPI spreadsheet.API
}

func NewAppendConfirmedBookingToTracker(spreadsheetAPI spreadsheet.API) *AppendConfirmedBookingToTracker {
	return &AppendConfirmedBookingToTracker{spreadsheetAPI: spreadsheetAPI}
}

type AppendConfirmedBookingToTrackerInput struct {
	TicketID      string
	CustomerEmail string
	Price         entity.Money
}

func (a *AppendConfirmedBookingToTracker) Execute(ctx context.Context, in AppendConfirmedBookingToTrackerInput) error {
	err := a.spreadsheetAPI.AppendRow(
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
