package handler

import (
	"net/http"
	"tickets/internal/domain/usecase/show"
	"time"

	"github.com/labstack/echo/v4"
)

type ShowHandler struct {
	createShowUseCase *show.Create
}

func NewShowHandler(createShowUseCase *show.Create) *ShowHandler {
	return &ShowHandler{createShowUseCase: createShowUseCase}
}

type createShowRequest struct {
	DeadNationID    string    `json:"dead_nation_id"`
	NumberOfTickets int       `json:"number_of_tickets"`
	StartTime       time.Time `json:"start_time"`
	Title           string    `json:"title"`
	Venue           string    `json:"venue"`
}

type createShowResponse struct {
	ID string `json:"id"`
}

func (h *ShowHandler) CreateShow(c echo.Context) error {
	var req createShowRequest
	err := c.Bind(&req)
	if err != nil {
		return err
	}

	id, err := h.createShowUseCase.Execute(c.Request().Context(), show.CreateInput{
		DeadNationID:    req.DeadNationID,
		NumberOfTickets: req.NumberOfTickets,
		StartTime:       req.StartTime,
		Title:           req.Title,
		Venue:           req.Venue,
	})
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, createShowResponse{ID: id})
}
