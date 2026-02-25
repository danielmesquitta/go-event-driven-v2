package handler

import (
	"net/http"
	"tickets/internal/domain/usecase/booking"

	"github.com/labstack/echo/v4"
)

type BookingHandler struct {
	createBookingUseCase *booking.Create
}

func NewBookingHandler(createBookingUseCase *booking.Create) *BookingHandler {
	return &BookingHandler{createBookingUseCase: createBookingUseCase}
}

type createBookingRequest struct {
	ShowID          string `json:"show_id"`
	NumberOfTickets int    `json:"number_of_tickets"`
	CustomerEmail   string `json:"customer_email"`
}

type createBookingResponse struct {
	BookingID string `json:"booking_id"`
}

func (h *BookingHandler) CreateBooking(c echo.Context) error {
	var req createBookingRequest
	err := c.Bind(&req)
	if err != nil {
		return err
	}

	id, err := h.createBookingUseCase.Execute(c.Request().Context(), booking.CreateInput{
		ShowID:          req.ShowID,
		NumberOfTickets: req.NumberOfTickets,
		CustomerEmail:   req.CustomerEmail,
	})
	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, createBookingResponse{BookingID: id})
}
