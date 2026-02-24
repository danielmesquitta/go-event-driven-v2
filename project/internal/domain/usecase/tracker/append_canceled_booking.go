package tracker

import (
	"context"
	"tickets/internal/domain/entity"
	"tickets/internal/pkg/validator"
	"tickets/internal/provider/spreadsheet"
)

type AppendCanceledBooking struct {
	spreadsheetAPI spreadsheet.API
}

func NewAppendCanceledBooking(spreadsheetAPI spreadsheet.API) *AppendCanceledBooking {
	return &AppendCanceledBooking{spreadsheetAPI: spreadsheetAPI}
}

type AppendCanceledBookingInput struct {
	TicketID      string       `json:"ticket_id" validate:"required"`
	CustomerEmail string       `json:"customer_email" validate:"required"`
	Price         entity.Money `json:"price" validate:"required"`
}

func (a *AppendCanceledBooking) Execute(ctx context.Context, in AppendCanceledBookingInput) error {
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
