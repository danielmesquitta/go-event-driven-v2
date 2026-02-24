package show

import (
	"context"
	"tickets/internal/domain/entity"
	"tickets/internal/pkg/validator"
	"tickets/internal/provider/repo"
	"time"

	"github.com/google/uuid"
)

type Create struct {
	showRepo repo.ShowRepo
}

func NewCreate(showRepo repo.ShowRepo) *Create {
	return &Create{showRepo: showRepo}
}

type CreateInput struct {
	DeadNationID    string    `json:"dead_nation_id" validate:"required"`
	NumberOfTickets int       `json:"number_of_tickets" validate:"required"`
	StartTime       time.Time `json:"start_time" validate:"required"`
	Title           string    `json:"title" validate:"required"`
	Venue           string    `json:"venue" validate:"required"`
}

func (c *Create) Execute(ctx context.Context, in CreateInput) (id string, err error) {
	err = validator.Validate(ctx, in)
	if err != nil {
		return "", err
	}

	id = uuid.NewString()
	err = c.showRepo.Create(ctx, &entity.Show{
		ID:              id,
		DeadNationID:    in.DeadNationID,
		NumberOfTickets: in.NumberOfTickets,
		StartTime:       in.StartTime,
		Title:           in.Title,
		Venue:           in.Venue,
	})
	if err != nil {
		return "", err
	}

	return id, nil
}
