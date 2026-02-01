// Package tracing provides OpenTelemetry wrapper: context propagation and minimal instrumentation.
// Uses noop by default; user can configure a real TracerProvider.
package tracing

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

// Tracer returns the tracer from the global provider (or noop if not configured).
func Tracer(name string) trace.Tracer {
	return otel.Tracer(name)
}

// StartSpan starts a span in the context and returns the context and span.
func StartSpan(ctx context.Context, tracerName, spanName string) (context.Context, trace.Span) {
	return otel.Tracer(tracerName).Start(ctx, spanName)
}

// SpanFromContext returns the span from the context, if any.
func SpanFromContext(ctx context.Context) trace.Span {
	return trace.SpanFromContext(ctx)
}
