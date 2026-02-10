// Package queue: instrumented consumer with OpenTelemetry and Prometheus metrics.

package queue

import (
	"context"
	"time"

	"github.com/cosmos-toolkit/pkgs/pkg/metrics"
	"github.com/cosmos-toolkit/pkgs/pkg/tracing"
	"go.opentelemetry.io/otel/attribute"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel/codes"
)

// InstrumentedConsumer wraps a Consumer and emits metrics and spans per message.
type InstrumentedConsumer struct {
	Consumer
	tracerName  string
	metricName  string
	msgTotal    prometheus.Counter
	msgErrors   prometheus.Counter
	msgLatency  prometheus.Observer
}

// InstrumentedConsumerConfig configures the instrumented consumer.
type InstrumentedConsumerConfig struct {
	TracerName string
	MetricName string
}

// DefaultInstrumentedConsumerConfig returns default instrumentation config.
func DefaultInstrumentedConsumerConfig() InstrumentedConsumerConfig {
	return InstrumentedConsumerConfig{
		TracerName: "queue",
		MetricName: "queue_messages",
	}
}

// NewInstrumentedConsumer wraps a Consumer with tracing and metrics.
func NewInstrumentedConsumer(c Consumer, cfg InstrumentedConsumerConfig) *InstrumentedConsumer {
	msgTotal := metrics.Counter(cfg.MetricName+"_total", "Total messages consumed")
	msgErrors := metrics.Counter(cfg.MetricName+"_errors_total", "Total message processing errors")
	msgLatency := metrics.Histogram(cfg.MetricName+"_duration_seconds", "Message processing latency in seconds", nil)

	return &InstrumentedConsumer{
		Consumer:   c,
		tracerName: cfg.TracerName,
		metricName: cfg.MetricName,
		msgTotal:   msgTotal,
		msgErrors:  msgErrors,
		msgLatency: msgLatency,
	}
}

// Consume runs the wrapped consumer with instrumented handler.
func (ic *InstrumentedConsumer) Consume(ctx context.Context, topic string, handler func(ctx context.Context, m *Message) error) error {
	return ic.Consumer.Consume(ctx, topic, func(ctx context.Context, m *Message) error {
		start := time.Now()
		ctx, span := tracing.StartSpan(ctx, ic.tracerName, "consume")
		span.SetAttributes(
			attribute.String("topic", topic),
			attribute.String("message_id", m.ID),
		)
		defer span.End()

		err := handler(ctx, m)
		ic.msgTotal.Inc()
		if err != nil {
			span.SetStatus(codes.Error, err.Error())
			span.SetAttributes(attribute.Bool("error", true))
			ic.msgErrors.Inc()
		}
		ic.msgLatency.Observe(time.Since(start).Seconds())
		return err
	})
}
