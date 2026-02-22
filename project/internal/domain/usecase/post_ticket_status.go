package usecase

import (
	"context"
	"fmt"
	"tickets/internal/app/pubsub/event"
	"tickets/internal/domain/entity"
	"tickets/internal/domain/errs"
	"tickets/internal/pkg/ctxval"
	"tickets/internal/pkg/validator"
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
	Tickets []entity.Ticket `json:"tickets"         validate:"required,dive"`
}

func (p *PostTicketStatus) Execute(ctx context.Context, in PostTicketStatusInput) error {
	err := validator.Validate(ctx, in)
	if err != nil {
		return err
	}

	g, gCtx := errgroup.WithContext(ctx)

	for _, ticket := range in.Tickets {
		g.Go(func() error {
			return p.postTicketStatus(gCtx, postTicketStatusInput{
				Ticket: ticket,
			})
		})
	}

	return g.Wait()
}

type postTicketStatusInput struct {
	Ticket entity.Ticket
}

func (p *PostTicketStatus) postTicketStatus(ctx context.Context, in postTicketStatusInput) error {
	idempotencyKey := ctxval.GetIdempotencyKey(ctx) + "-" + in.Ticket.ID
	ctx = ctxval.WithIdempotencyKey(ctx, idempotencyKey)

	var e event.Event
	switch in.Ticket.Status {
	case entity.TicketStatusConfirmed:
		e = event.NewTicketBookingConfirmed(
			ctx,
			in.Ticket.ID,
			in.Ticket.CustomerEmail,
			in.Ticket.Price,
		)

	case entity.TicketStatusCanceled:
		e = event.NewTicketBookingCanceled(
			ctx,
			in.Ticket.ID,
			in.Ticket.CustomerEmail,
			in.Ticket.Price,
		)

	default:
		msgs := map[string]string{
			"status": fmt.Sprintf("status must be confirmed or canceled, got: %s", in.Ticket.Status),
		}
		return errs.ErrInvalidFormat.New(errs.WithMetadata(errs.MetadataErrorsKey, msgs))
	}

	if err := p.eventBus.Publish(ctx, e); err != nil {
		return err
	}

	return nil
}
