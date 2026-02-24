package router

import (
	"tickets/internal/app/http/handler"
	"tickets/internal/app/http/middleware"
	"tickets/internal/domain/usecase/ticket"

	libHTTP "github.com/ThreeDotsLabs/go-event-driven/v2/common/http"
	"github.com/labstack/echo/v4"
)

func New(
	postTicketStatusUseCase *ticket.PostStatus,
	listTicketsUseCase *ticket.List,
) *echo.Echo {
	e := libHTTP.NewEcho()

	ticketHandler := handler.NewTicketHandler(
		postTicketStatusUseCase,
		listTicketsUseCase,
	)

	healthHandler := handler.NewHealthHandler()

	e.Use(
		middleware.CorrelationID,
		middleware.IdempotencyKey,
	)

	e.GET("/health", healthHandler.Handle)

	e.POST("/tickets-status", ticketHandler.PostTicketsStatus)
	e.GET("/tickets", ticketHandler.ListTickets)

	return e
}
