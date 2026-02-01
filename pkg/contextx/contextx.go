// Package contextx provides helpers for context.Context: default timeout,
// metadata (tenant, trace, user) and cascade cancellation.
package contextx

import (
	"context"
	"time"
)

type contextKey int

const (
	keyTraceID contextKey = iota
	keyRequestID
	keyUserID
	keyTenantID
)

// DefaultTimeout is the default timeout for operations (e.g. HTTP, DB).
const DefaultTimeout = 30 * time.Second

// WithTimeout returns a context with timeout; uses DefaultTimeout if d <= 0.
func WithTimeout(ctx context.Context, d time.Duration) (context.Context, context.CancelFunc) {
	if d <= 0 {
		d = DefaultTimeout
	}
	return context.WithTimeout(ctx, d)
}

// WithTraceID sets trace_id in the context.
func WithTraceID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, keyTraceID, id)
}

// TraceID returns the trace_id from the context.
func TraceID(ctx context.Context) string {
	v := ctx.Value(keyTraceID)
	if s, ok := v.(string); ok {
		return s
	}
	return ""
}

// WithRequestID sets request_id in the context.
func WithRequestID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, keyRequestID, id)
}

// RequestID returns the request_id from the context.
func RequestID(ctx context.Context) string {
	v := ctx.Value(keyRequestID)
	if s, ok := v.(string); ok {
		return s
	}
	return ""
}

// WithUserID sets user_id in the context.
func WithUserID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, keyUserID, id)
}

// UserID returns the user_id from the context.
func UserID(ctx context.Context) string {
	v := ctx.Value(keyUserID)
	if s, ok := v.(string); ok {
		return s
	}
	return ""
}

// WithTenant sets tenant_id in the context.
func WithTenant(ctx context.Context, tenantID string) context.Context {
	return context.WithValue(ctx, keyTenantID, tenantID)
}

// Tenant returns the tenant_id from the context.
func Tenant(ctx context.Context) string {
	v := ctx.Value(keyTenantID)
	if s, ok := v.(string); ok {
		return s
	}
	return ""
}
