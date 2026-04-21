package pg

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"tickets/internal/domain/entity"
	"tickets/internal/provider/db"
	"tickets/internal/provider/repo"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
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

func (r *OpsBooking) GetByTicketID(ctx context.Context, ticketID string) (*entity.OpsBooking, error) {
	var payload []byte
	err := r.db.WithTx(ctx).GetContext(ctx, &payload, `
		SELECT payload FROM read_model_ops_bookings
		WHERE payload->'tickets' ? $1
		LIMIT 1
	`, ticketID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get ops booking for ticket %s: %w", ticketID, err)
	}

	op, err := unmarshalOpsBooking(payload)
	if err != nil {
		return nil, err
	}
	return &op, nil
}

func (r *OpsBooking) GetByBookingID(ctx context.Context, bookingID string) (*entity.OpsBooking, error) {
	var payload []byte
	err := r.db.WithTx(ctx).GetContext(ctx, &payload, `
		SELECT payload FROM read_model_ops_bookings
		WHERE booking_id = $1
		LIMIT 1
	`, bookingID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get ops booking for booking %s: %w", bookingID, err)
	}

	op, err := unmarshalOpsBooking(payload)
	if err != nil {
		return nil, err
	}
	return &op, nil
}

func (r *OpsBooking) UpdateByBookingID(
	ctx context.Context,
	bookingID string,
	updateFn func(entity.OpsBooking) (entity.OpsBooking, error),
) error {
	return r.updateInTx(ctx, func(ctx context.Context, tx *sqlx.Tx) error {
		op, err := r.findByBookingID(ctx, tx, bookingID)
		if errors.Is(err, sql.ErrNoRows) {
			// Events may arrive out of order; retry until the read model exists.
			return fmt.Errorf("read model for booking %s does not exist yet", bookingID)
		} else if err != nil {
			return fmt.Errorf("could not find read model: %w", err)
		}

		updated, err := updateFn(op)
		if err != nil {
			return err
		}

		return r.updateReadModel(ctx, tx, updated)
	})
}

func (r *OpsBooking) UpdateByTicketID(
	ctx context.Context,
	ticketID string,
	updateFn func(entity.OpsBooking) (entity.OpsBooking, error),
) error {
	return r.updateInTx(ctx, func(ctx context.Context, tx *sqlx.Tx) error {
		op, err := r.findByTicketID(ctx, tx, ticketID)
		if errors.Is(err, sql.ErrNoRows) {
			// Events may arrive out of order; retry until the ticket is registered.
			return fmt.Errorf("read model for ticket %s does not exist yet", ticketID)
		} else if err != nil {
			return fmt.Errorf("could not find read model: %w", err)
		}

		updated, err := updateFn(op)
		if err != nil {
			return err
		}

		return r.updateReadModel(ctx, tx, updated)
	})
}

func (r *OpsBooking) updateInTx(
	ctx context.Context,
	fn func(ctx context.Context, tx *sqlx.Tx) error,
) (err error) {
	tx, err := r.db.BeginTxx(ctx, &sql.TxOptions{Isolation: sql.LevelRepeatableRead})
	if err != nil {
		return fmt.Errorf("could not begin transaction: %w", err)
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
			return
		}
		if commitErr := tx.Commit(); commitErr != nil {
			err = fmt.Errorf("could not commit transaction: %w", commitErr)
		}
	}()

	return fn(ctx, tx)
}

func (r *OpsBooking) updateReadModel(ctx context.Context, tx *sqlx.Tx, op entity.OpsBooking) error {
	op.LastUpdate = time.Now().UTC()

	payload, err := json.Marshal(op)
	if err != nil {
		return fmt.Errorf("could not marshal ops booking: %w", err)
	}

	_, err = tx.ExecContext(ctx, `
		INSERT INTO read_model_ops_bookings (booking_id, payload)
		VALUES ($1, $2)
		ON CONFLICT (booking_id) DO UPDATE SET payload = EXCLUDED.payload
	`, op.BookingID, payload)
	if err != nil {
		return fmt.Errorf("could not update ops booking: %w", err)
	}

	return nil
}

func (r *OpsBooking) findByBookingID(ctx context.Context, tx *sqlx.Tx, bookingID string) (entity.OpsBooking, error) {
	var payload []byte
	err := tx.GetContext(ctx, &payload, `
		SELECT payload FROM read_model_ops_bookings
		WHERE booking_id = $1
	`, bookingID)
	if err != nil {
		return entity.OpsBooking{}, err
	}
	return unmarshalOpsBooking(payload)
}

func (r *OpsBooking) findByTicketID(ctx context.Context, tx *sqlx.Tx, ticketID string) (entity.OpsBooking, error) {
	var payload []byte
	err := tx.GetContext(ctx, &payload, `
		SELECT payload FROM read_model_ops_bookings
		WHERE payload->'tickets' ? $1
		LIMIT 1
	`, ticketID)
	if err != nil {
		return entity.OpsBooking{}, err
	}
	return unmarshalOpsBooking(payload)
}

func unmarshalOpsBooking(payload []byte) (entity.OpsBooking, error) {
	var op entity.OpsBooking
	if err := json.Unmarshal(payload, &op); err != nil {
		return entity.OpsBooking{}, fmt.Errorf("could not unmarshal ops booking: %w", err)
	}
	if op.Tickets == nil {
		op.Tickets = map[string]entity.OpsTicket{}
	}
	return op, nil
}

var _ repo.OpsBookingRepo = &OpsBooking{}
