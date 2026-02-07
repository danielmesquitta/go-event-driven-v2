package log

import (
	"context"
	"log/slog"
	"os"
)

type ctxKey int

const (
	loggerKey ctxKey = iota
)

// FromContext returns the logger from the context.
func FromContext(ctx context.Context) *slog.Logger {
	log, ok := ctx.Value(loggerKey).(*slog.Logger)
	if ok {
		return log
	}

	return slog.New(slog.NewTextHandler(os.Stderr, nil))
}

// ToContext adds the logger to the context.
func ToContext(ctx context.Context, logger *slog.Logger) context.Context {
	return context.WithValue(ctx, loggerKey, logger)
}
