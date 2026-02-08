package handler

import (
	"context"
	"tickets/internal/app/pubsub/event"
	"tickets/internal/provider/spreadsheet"
)

type AppendConfirmedBookingToTracker struct {
	spreadsheetAPI spreadsheet.API
}

func NewAppendConfirmedBookingToTracker(spreadsheetAPI spreadsheet.API) *AppendConfirmedBookingToTracker {
	return &AppendConfirmedBookingToTracker{spreadsheetAPI: spreadsheetAPI}
}

func (a *AppendConfirmedBookingToTracker) Handle(ctx context.Context, event *event.TicketBookingConfirmed) error {
	err := a.spreadsheetAPI.AppendRow(
		ctx,
		"tickets-to-print",
		[]string{
			event.TicketID,
			event.CustomerEmail,
			event.Price.Amount,
			event.Price.Currency,
		},
	)
	if err != nil {
		return err
	}
	return nil
}
