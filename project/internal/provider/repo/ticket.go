package repo

import (
	"context"
	"tickets/internal/domain/entity"
)

type TicketRepo interface {
	Create(ctx context.Context, ticket *entity.Ticket) error
	Get(ctx context.Context, id string) (*entity.Ticket, error)
	List(ctx context.Context) ([]entity.Ticket, error)
	Delete(ctx context.Context, id string) error
}
