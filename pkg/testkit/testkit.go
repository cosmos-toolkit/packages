// Package testkit provides builders, fakes and helpers for tests
// (nop logger, context, etc.).
package testkit

import (
	"context"
	"io"
	"log/slog"

	"github.com/cosmos-toolkit/pkgs/pkg/contextx"
)

// NopLogger returns a slog.Logger that discards all output (useful in tests).
func NopLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(io.Discard, &slog.HandlerOptions{Level: slog.LevelError}))
}

// ContextWithIDs returns a context with trace_id, request_id and optional user_id set.
func ContextWithIDs(traceID, requestID, userID string) context.Context {
	ctx := context.Background()
	if traceID != "" {
		ctx = contextx.WithTraceID(ctx, traceID)
	}
	if requestID != "" {
		ctx = contextx.WithRequestID(ctx, requestID)
	}
	if userID != "" {
		ctx = contextx.WithUserID(ctx, userID)
	}
	return ctx
}
