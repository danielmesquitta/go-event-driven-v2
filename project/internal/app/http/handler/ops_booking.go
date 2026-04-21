package handler

import (
	"net/http"
	"time"

	"tickets/internal/domain/entity"
	"tickets/internal/domain/errs"
	"tickets/internal/domain/usecase/opsbooking"

	"github.com/labstack/echo/v4"
)

const receiptIssueDateParam = "receipt_issue_date"
const receiptIssueDateLayout = "2006-01-02"

type OpsBookingHandler struct {
	listOpsBookingsUseCase *opsbooking.List
	getOpsBookingUseCase   *opsbooking.Get
}

func NewOpsBookingHandler(
	listOpsBookingsUseCase *opsbooking.List,
	getOpsBookingUseCase *opsbooking.Get,
) *OpsBookingHandler {
	return &OpsBookingHandler{
		listOpsBookingsUseCase: listOpsBookingsUseCase,
		getOpsBookingUseCase:   getOpsBookingUseCase,
	}
}

func (h *OpsBookingHandler) ListOpsBookings(c echo.Context) error {
	input := opsbooking.ListInput{}

	if raw := c.QueryParam(receiptIssueDateParam); raw != "" {
		date, err := time.Parse(receiptIssueDateLayout, raw)
		if err != nil {
			return errs.ErrInvalidFormat.New(errs.WithMetadata(
				errs.MetadataErrorsKey,
				map[string]string{
					receiptIssueDateParam: "must be a date in YYYY-MM-DD format",
				},
			))
		}
		input.ReceiptIssueDate = &date
	}

	bookings, err := h.listOpsBookingsUseCase.Execute(c.Request().Context(), input)
	if err != nil {
		return err
	}

	if bookings == nil {
		bookings = []entity.OpsBooking{}
	}

	return c.JSON(http.StatusOK, bookings)
}

func (h *OpsBookingHandler) GetOpsBooking(c echo.Context) error {
	bookingID := c.Param("id")

	booking, err := h.getOpsBookingUseCase.Execute(c.Request().Context(), opsbooking.GetInput{
		BookingID: bookingID,
	})
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, booking)
}
