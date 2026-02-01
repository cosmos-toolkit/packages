// Package logger: context helpers re-export contextx for convenience.
// To set trace_id/request_id in context, use contextx.WithTraceID etc.
package logger

import (
	"context"

	"github.com/cosmos-toolkit/pkgs/pkg/contextx"
)

// TraceID returns the trace_id from context (delegates to contextx).
func TraceID(ctx context.Context) string { return contextx.TraceID(ctx) }

// RequestID returns the request_id from context (delegates to contextx).
func RequestID(ctx context.Context) string { return contextx.RequestID(ctx) }

// UserID returns the user_id from context (delegates to contextx).
func UserID(ctx context.Context) string { return contextx.UserID(ctx) }
