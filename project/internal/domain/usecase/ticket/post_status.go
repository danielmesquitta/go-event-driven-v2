package ticket

import (
	"context"
	"fmt"
	"tickets/internal/app/pubsub/message/event"
	"tickets/internal/domain/entity"
	"tickets/internal/domain/errs"
	"tickets/internal/pkg/ctxval"
	"tickets/internal/pkg/validator"
	"tickets/internal/provider/eventbus"

	"golang.org/x/sync/errgroup"
)

type PostStatus struct {
	eventBus eventbus.EventBus
}

func NewPostStatus(
	eventBus eventbus.EventBus,
) *PostStatus {
	return &PostStatus{
		eventBus: eventBus,
	}
}

type PostStatusInput struct {
	Tickets []entity.Ticket `json:"tickets"         validate:"required,dive"`
}

func (p *PostStatus) Execute(ctx context.Context, in PostStatusInput) error {
	err := validator.Validate(ctx, in)
	if err != nil {
		return err
	}

	g, gCtx := errgroup.WithContext(ctx)

	for _, t := range in.Tickets {
		g.Go(func() error {
			return p.postTicketStatus(gCtx, postTicketStatusInput{
				Ticket: t,
			})
		})
	}

	return g.Wait()
}

type postTicketStatusInput struct {
	Ticket entity.Ticket
}

func (p *PostStatus) postTicketStatus(ctx context.Context, in postTicketStatusInput) error {
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
