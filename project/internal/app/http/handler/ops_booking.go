package handler

import (
	"net/http"
	"tickets/internal/domain/entity"
	"tickets/internal/domain/usecase/opsbooking"

	"github.com/labstack/echo/v4"
)

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
	bookings, err := h.listOpsBookingsUseCase.Execute(c.Request().Context())
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
