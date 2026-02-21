package usecase

import (
	"context"
	"tickets/internal/domain/entity"
	"tickets/internal/provider/repo"
)

type ListTickets struct {
	ticketRepo repo.TicketRepo
}

func NewListTickets(ticketRepo repo.TicketRepo) *ListTickets {
	return &ListTickets{ticketRepo: ticketRepo}
}

func (l *ListTickets) Execute(ctx context.Context) ([]entity.Ticket, error) {
	return l.ticketRepo.List(ctx)
}
