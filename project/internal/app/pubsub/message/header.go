package message

import (
	"context"
	"tickets/internal/pkg/ctxval"
	"time"

	"github.com/google/uuid"
)

type Header struct {
	ID             string    `json:"id"`
	PublishedAt    time.Time `json:"published_at"`
	IdempotencyKey string    `json:"idempotency_key"`
}

func NewHeader(ctx context.Context) Header {
	return Header{
		ID:             uuid.NewString(),
		PublishedAt:    time.Now(),
		IdempotencyKey: ctxval.GetIdempotencyKey(ctx),
	}
}
