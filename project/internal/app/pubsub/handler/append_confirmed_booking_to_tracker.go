package handler

import (
	"encoding/json"
	"tickets/internal/app/pubsub/event"
	"tickets/internal/domain/errs"
	"tickets/internal/provider/spreadsheet"

	"github.com/ThreeDotsLabs/watermill/message"
)

type AppendConfirmedBookingToTracker struct {
	spreadsheetAPI spreadsheet.API
}

func NewAppendConfirmedBookingToTracker(spreadsheetAPI spreadsheet.API) *AppendConfirmedBookingToTracker {
	return &AppendConfirmedBookingToTracker{spreadsheetAPI: spreadsheetAPI}
}

func (a *AppendConfirmedBookingToTracker) Handle(msg *message.Message) error {
	var e event.TicketBookingConfirmed
	err := json.Unmarshal(msg.Payload, &e)
	if err != nil {
		return errs.ErrInvalidFormat.New(errs.WithMetadata("payload", string(msg.Payload)))
	}

	e.SetDefaults()

	err = a.spreadsheetAPI.AppendRow(
		msg.Context(),
		"tickets-to-print",
		[]string{
			e.TicketID,
			e.CustomerEmail,
			e.Price.Amount,
			e.Price.Currency,
		},
	)
	if err != nil {
		return err
	}
	return nil
}
