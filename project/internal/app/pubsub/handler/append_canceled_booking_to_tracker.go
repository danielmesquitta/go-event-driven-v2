package handler

import (
	"context"
	"tickets/internal/app/pubsub/event"
	"tickets/internal/provider/spreadsheet"
)

type AppendCanceledBookingToTracker struct {
	spreadsheetAPI spreadsheet.API
}

func NewAppendCanceledBookingToTracker(spreadsheetAPI spreadsheet.API) *AppendCanceledBookingToTracker {
	return &AppendCanceledBookingToTracker{spreadsheetAPI: spreadsheetAPI}
}

func (a *AppendCanceledBookingToTracker) Handle(ctx context.Context, event *event.TicketBookingCanceled) error {
	err := a.spreadsheetAPI.AppendRow(
		ctx,
		"tickets-to-refund",
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
