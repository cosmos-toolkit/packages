// Package logger: helpers de contexto re-exportam contextx para conveniência.
// Para definir trace_id/request_id no contexto, use contextx.WithTraceID etc.
package logger

import (
	"context"

	"github.com/cosmos-toolkit/pkgs/pkg/contextx"
)

// TraceID retorna o trace_id do contexto (delega para contextx).
func TraceID(ctx context.Context) string { return contextx.TraceID(ctx) }

// RequestID retorna o request_id do contexto (delega para contextx).
func RequestID(ctx context.Context) string { return contextx.RequestID(ctx) }

// UserID retorna o user_id do contexto (delega para contextx).
func UserID(ctx context.Context) string { return contextx.UserID(ctx) }
