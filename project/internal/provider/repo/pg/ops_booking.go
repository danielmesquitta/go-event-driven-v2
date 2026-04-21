package pg

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"tickets/internal/app/pubsub/message/event"
	"tickets/internal/domain/entity"
	"tickets/internal/provider/db"
	"tickets/internal/provider/repo"

	"github.com/google/uuid"
)

type OpsBooking struct {
	db *db.DB
}

func NewOpsBooking(db *db.DB) *OpsBooking {
	return &OpsBooking{db: db}
}

func (r *OpsBooking) Create(ctx context.Context, in repo.CreateOpsBookingInput) error {
	bookingID, err := uuid.Parse(in.BookingID)
	if err != nil {
		return fmt.Errorf("could not parse booking ID: %w", err)
	}

	op := entity.OpsBooking{
		BookingID:  bookingID,
		BookedAt:   in.BookedAt,
		Tickets:    map[string]entity.OpsTicket{},
		LastUpdate: time.Now().UTC(),
	}

	return r.createReadModel(ctx, op)
}

func (r *OpsBooking) Update(ctx context.Context, op entity.OpsBooking) error {
	op.LastUpdate = time.Now().UTC()

	updatedPayload, err := json.Marshal(op)
	if err != nil {
		return fmt.Errorf("could not marshal ops booking: %w", err)
	}

	_, err = r.db.WithTx(ctx).ExecContext(ctx, `
		UPDATE read_model_ops_bookings SET payload = $1 WHERE booking_id = $2
	`, updatedPayload, op.BookingID)
	if err != nil {
		return fmt.Errorf("could not update ops booking: %w", err)
	}

	return nil
}

func (r *OpsBooking) GetByTicketID(ctx context.Context, ticketID string) (*entity.OpsBooking, error) {
	var row struct {
		BookingID string `db:"booking_id"`
		Payload   []byte `db:"payload"`
	}
	err := r.db.WithTx(ctx).GetContext(ctx, &row, `
		SELECT booking_id, payload FROM read_model_ops_bookings
		WHERE payload->'tickets' ? $1
		LIMIT 1
	`, ticketID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get ops booking for ticket %s: %w", ticketID, err)
	}

	var op entity.OpsBooking
	if err = json.Unmarshal(row.Payload, &op); err != nil {
		return nil, fmt.Errorf("could not unmarshal ops booking: %w", err)
	}

	return &op, nil
}

func (r *OpsBooking) GetByBookingID(ctx context.Context, bookingID string) (*entity.OpsBooking, error) {
	var row struct {
		BookingID string `db:"booking_id"`
		Payload   []byte `db:"payload"`
	}
	err := r.db.WithTx(ctx).GetContext(ctx, &row, `
		SELECT booking_id, payload FROM read_model_ops_bookings
		WHERE booking_id = $1
		LIMIT 1
	`, bookingID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get ops booking for booking %s: %w", bookingID, err)
	}

	var op entity.OpsBooking
	if err = json.Unmarshal(row.Payload, &op); err != nil {
		return nil, fmt.Errorf("could not unmarshal ops booking: %w", err)
	}

	return &op, nil
}

// OnTicketBookingConfirmed adds or updates a ticket within the booking ops booking.
func (r *OpsBooking) OnTicketBookingConfirmed(ctx context.Context, e *event.TicketBookingConfirmed) error {
	return r.updateReadModelByBookingID(
		ctx,
		e.BookingID,
		func(op entity.OpsBooking) (entity.OpsBooking, error) {
			ticket, ok := op.Tickets[e.TicketID]
			if !ok {
				slog.InfoContext(ctx, "Creating ticket ops booking for ticket "+e.TicketID, "ticket_id", e.TicketID)
			}

			ticket.PriceAmount = e.Price.Amount
			ticket.PriceCurrency = e.Price.Currency
			ticket.CustomerEmail = e.CustomerEmail
			ticket.Status = "confirmed"
			op.Tickets[e.TicketID] = ticket

			return op, nil
		},
	)
}

// OnTicketRefunded sets the ticket status to "refunded".
func (r *OpsBooking) OnTicketRefunded(ctx context.Context, e *event.TicketRefunded) error {
	return r.updateReadModelByTicketID(
		ctx,
		e.TicketID,
		func(op entity.OpsBooking) (entity.OpsBooking, error) {
			ticket, ok := op.Tickets[e.TicketID]
			if !ok {
				return op, fmt.Errorf("ticket %s not found in ops booking", e.TicketID)
			}

			ticket.Status = "refunded"
			op.Tickets[e.TicketID] = ticket

			return op, nil
		},
	)
}

// OnTicketPrinted records printing details on the ticket.
func (r *OpsBooking) OnTicketPrinted(ctx context.Context, e *event.TicketPrinted) error {
	return r.updateReadModelByTicketID(
		ctx,
		e.TicketID,
		func(op entity.OpsBooking) (entity.OpsBooking, error) {
			ticket, ok := op.Tickets[e.TicketID]
			if !ok {
				return op, fmt.Errorf("ticket %s not found in ops booking", e.TicketID)
			}

			ticket.PrintedAt = e.Header.PublishedAt
			ticket.PrintedFileName = e.FileName
			op.Tickets[e.TicketID] = ticket

			return op, nil
		},
	)
}

// OnTicketReceiptIssued records receipt details on the ticket.
func (r *OpsBooking) OnTicketReceiptIssued(ctx context.Context, e *event.TicketReceiptIssued) error {
	return r.updateReadModelByTicketID(
		ctx,
		e.TicketID,
		func(op entity.OpsBooking) (entity.OpsBooking, error) {
			ticket, ok := op.Tickets[e.TicketID]
			if !ok {
				return op, fmt.Errorf("ticket %s not found in ops booking", e.TicketID)
			}

			ticket.ReceiptIssuedAt = e.IssuedAt
			ticket.ReceiptNumber = e.ReceiptNumber
			op.Tickets[e.TicketID] = ticket

			return op, nil
		},
	)
}

func (r *OpsBooking) createReadModel(ctx context.Context, op entity.OpsBooking) error {
	payload, err := json.Marshal(op)
	if err != nil {
		return fmt.Errorf("could not marshal ops booking: %w", err)
	}

	_, err = r.db.WithTx(ctx).ExecContext(ctx, `
		INSERT INTO read_model_ops_bookings (booking_id, payload)
		VALUES ($1, $2)
		ON CONFLICT (booking_id) DO NOTHING
	`, op.BookingID, payload)
	if err != nil {
		return fmt.Errorf("could not create ops booking: %w", err)
	}

	return nil
}

func (r *OpsBooking) updateReadModelByBookingID(
	ctx context.Context,
	bookingID string,
	updateFn func(entity.OpsBooking) (entity.OpsBooking, error),
) error {
	var payload []byte
	err := r.db.WithTx(ctx).GetContext(ctx, &payload, `
		SELECT payload FROM read_model_ops_bookings WHERE booking_id = $1 FOR UPDATE
	`, bookingID)
	if err != nil {
		return fmt.Errorf("could not get ops booking for booking %s: %w", bookingID, err)
	}

	var op entity.OpsBooking
	if err = json.Unmarshal(payload, &op); err != nil {
		return fmt.Errorf("could not unmarshal ops booking: %w", err)
	}

	op, err = updateFn(op)
	if err != nil {
		return err
	}
	op.LastUpdate = time.Now().UTC()

	updatedPayload, err := json.Marshal(op)
	if err != nil {
		return fmt.Errorf("could not marshal ops booking: %w", err)
	}

	_, err = r.db.WithTx(ctx).ExecContext(ctx, `
		UPDATE read_model_ops_bookings SET payload = $1 WHERE booking_id = $2
	`, updatedPayload, bookingID)
	if err != nil {
		return fmt.Errorf("could not update ops booking: %w", err)
	}

	return nil
}

func (r *OpsBooking) updateReadModelByTicketID(
	ctx context.Context,
	ticketID string,
	updateFn func(entity.OpsBooking) (entity.OpsBooking, error),
) error {
	var row struct {
		BookingID string `db:"booking_id"`
		Payload   []byte `db:"payload"`
	}
	err := r.db.WithTx(ctx).GetContext(ctx, &row, `
		SELECT booking_id, payload FROM read_model_ops_bookings
		WHERE payload->'tickets' ? $1
		FOR UPDATE
	`, ticketID)
	if err != nil {
		return fmt.Errorf("could not find ops booking for ticket %s: %w", ticketID, err)
	}

	var op entity.OpsBooking
	if err = json.Unmarshal(row.Payload, &op); err != nil {
		return fmt.Errorf("could not unmarshal ops booking: %w", err)
	}

	op, err = updateFn(op)
	if err != nil {
		return err
	}
	op.LastUpdate = time.Now().UTC()

	updatedPayload, err := json.Marshal(op)
	if err != nil {
		return fmt.Errorf("could not marshal ops booking: %w", err)
	}

	_, err = r.db.WithTx(ctx).ExecContext(ctx, `
		UPDATE read_model_ops_bookings SET payload = $1 WHERE booking_id = $2
	`, updatedPayload, row.BookingID)
	if err != nil {
		return fmt.Errorf("could not update ops booking: %w", err)
	}

	return nil
}

var _ repo.OpsBookingRepo = &OpsBooking{}
