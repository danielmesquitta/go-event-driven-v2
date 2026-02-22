package event

import (
	"context"
	"tickets/internal/pkg/ctxval"
	"time"

	"github.com/google/uuid"
)

type EventHeader struct {
	ID             string    `json:"id"`
	PublishedAt    time.Time `json:"published_at"`
	IdempotencyKey string    `json:"idempotency_key"`
}

func NewEventHeader(ctx context.Context) EventHeader {
	return EventHeader{
		ID:             uuid.NewString(),
		PublishedAt:    time.Now(),
		IdempotencyKey: ctxval.GetIdempotencyKey(ctx),
	}
}
