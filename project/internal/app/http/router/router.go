package router

import (
	"tickets/internal/app/http/handler"
	"tickets/internal/app/http/middleware"

	libHTTP "github.com/ThreeDotsLabs/go-event-driven/v2/common/http"
	"github.com/labstack/echo/v4"
)

func New(
	showHandler *handler.ShowHandler,
	ticketHandler *handler.TicketHandler,
	healthHandler *handler.HealthHandler,
	bookingHandler *handler.BookingHandler,
	opsBookingHandler *handler.OpsBookingHandler,
) *echo.Echo {
	e := libHTTP.NewEcho()

	e.Use(
		middleware.CorrelationID,
		middleware.IdempotencyKey,
		middleware.ErrorHandler,
	)

	e.GET("/health", healthHandler.Handle)

	e.POST("/shows", showHandler.CreateShow)

	e.POST("/tickets-status", ticketHandler.PostTicketsStatus)
	e.GET("/tickets", ticketHandler.ListTickets)
	e.PUT("/ticket-refund/:ticket_id", ticketHandler.TicketRefund)

	e.POST("/book-tickets", bookingHandler.CreateBooking)

	e.GET("/ops/bookings", opsBookingHandler.ListOpsBookings)
	e.GET("/ops/bookings/:id", opsBookingHandler.GetOpsBooking)

	return e
}
