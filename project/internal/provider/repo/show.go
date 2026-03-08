package repo

import (
	"context"
	"tickets/internal/domain/entity"
)

type ShowRepo interface {
	Create(ctx context.Context, show *entity.Show) error
	Get(ctx context.Context, id string) (*entity.Show, error)
	AvailableTickets(ctx context.Context, id string) (int, error)
}
