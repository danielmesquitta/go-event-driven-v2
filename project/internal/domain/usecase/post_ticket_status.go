package usecase

import (
	"context"
	"tickets/internal/app/pubsub/event"
	"tickets/internal/domain/entity"
	"tickets/internal/domain/errs"
	"tickets/internal/provider/eventbus"

	"golang.org/x/sync/errgroup"
)

type PostTicketStatus struct {
	eventBus eventbus.EventBus
}

func NewPostTicketStatus(
	eventBus eventbus.EventBus,
) *PostTicketStatus {
	return &PostTicketStatus{
		eventBus: eventBus,
	}
}

type PostTicketStatusInput struct {
	CorrelationID string          `json:"correlation_id"`
	Tickets       []entity.Ticket `json:"tickets"`
}

func (p *PostTicketStatus) Execute(ctx context.Context, in PostTicketStatusInput) error {
	g, gCtx := errgroup.WithContext(ctx)

	for _, ticket := range in.Tickets {
		g.Go(func() error {
			return p.postTicketStatus(gCtx, ticket, in.CorrelationID)
		})
	}

	return g.Wait()
}

func (p *PostTicketStatus) postTicketStatus(ctx context.Context, ticket entity.Ticket, correlationID string) error {
	var e event.Event
	switch ticket.Status {
	case entity.TicketStatusConfirmed:
		e = &event.TicketBookingConfirmed{
			Header:        event.NewEventHeader(),
			TicketID:      ticket.ID,
			CustomerEmail: ticket.CustomerEmail,
			Price:         ticket.Price,
		}

	case entity.TicketStatusCanceled:
		e = &event.TicketBookingCanceled{
			Header:        event.NewEventHeader(),
			TicketID:      ticket.ID,
			CustomerEmail: ticket.CustomerEmail,
			Price:         ticket.Price,
		}

	default:
		return errs.ErrInvalidFormat.New(errs.WithMetadata("status", ticket.Status))
	}

	if err := p.eventBus.Publish(ctx, e, eventbus.WithCorrelationID(correlationID)); err != nil {
		return err
	}

	return nil
}
