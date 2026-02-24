package ticket

import (
	"context"
	"tickets/internal/domain/entity"
	"tickets/internal/provider/repo"
)

type List struct {
	ticketRepo repo.TicketRepo
}

func NewList(ticketRepo repo.TicketRepo) *List {
	return &List{ticketRepo: ticketRepo}
}

func (l *List) Execute(ctx context.Context) ([]entity.Ticket, error) {
	return l.ticketRepo.List(ctx)
}
