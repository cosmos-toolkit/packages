// Package contextx fornece helpers para context.Context: timeout padrão,
// metadata (tenant, trace, user) e cancelamento em cascata.
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

// DefaultTimeout é o timeout padrão para operações (ex.: HTTP, DB).
const DefaultTimeout = 30 * time.Second

// WithTimeout retorna um contexto com timeout; usa DefaultTimeout se d <= 0.
func WithTimeout(ctx context.Context, d time.Duration) (context.Context, context.CancelFunc) {
	if d <= 0 {
		d = DefaultTimeout
	}
	return context.WithTimeout(ctx, d)
}

// WithTraceID coloca trace_id no contexto.
func WithTraceID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, keyTraceID, id)
}

// TraceID retorna o trace_id do contexto.
func TraceID(ctx context.Context) string {
	v := ctx.Value(keyTraceID)
	if s, ok := v.(string); ok {
		return s
	}
	return ""
}

// WithRequestID coloca request_id no contexto.
func WithRequestID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, keyRequestID, id)
}

// RequestID retorna o request_id do contexto.
func RequestID(ctx context.Context) string {
	v := ctx.Value(keyRequestID)
	if s, ok := v.(string); ok {
		return s
	}
	return ""
}

// WithUserID coloca user_id no contexto.
func WithUserID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, keyUserID, id)
}

// UserID retorna o user_id do contexto.
func UserID(ctx context.Context) string {
	v := ctx.Value(keyUserID)
	if s, ok := v.(string); ok {
		return s
	}
	return ""
}

// WithTenant coloca tenant_id no contexto.
func WithTenant(ctx context.Context, tenantID string) context.Context {
	return context.WithValue(ctx, keyTenantID, tenantID)
}

// Tenant retorna o tenant_id do contexto.
func Tenant(ctx context.Context) string {
	v := ctx.Value(keyTenantID)
	if s, ok := v.(string); ok {
		return s
	}
	return ""
}
