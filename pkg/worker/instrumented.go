// Package worker: instrumented pool with OpenTelemetry and Prometheus metrics.

package worker

import (
	"context"
	"time"

	"github.com/cosmos-toolkit/pkgs/pkg/metrics"
	"github.com/cosmos-toolkit/pkgs/pkg/tracing"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

// InstrumentedPool wraps a Pool and instruments each job with tracing and metrics.
type InstrumentedPool struct {
	*Pool
}

// InstrumentedPoolConfig configures the instrumented pool.
type InstrumentedPoolConfig struct {
	TracerName string
	JobName   string
}

// DefaultInstrumentedConfig returns default instrumentation config.
func DefaultInstrumentedConfig() InstrumentedPoolConfig {
	return InstrumentedPoolConfig{
		TracerName: "worker",
		JobName:   "worker_jobs",
	}
}

// NewInstrumentedPool creates a pool that instruments each job with spans and metrics.
func NewInstrumentedPool(cfg Config, jobsCh <-chan Job, inst InstrumentedPoolConfig) *InstrumentedPool {
	latency := metrics.Histogram(inst.JobName+"_duration_seconds", "Job execution latency in seconds", nil)
	total := metrics.Counter(inst.JobName+"_total", "Total jobs processed")
	errors := metrics.Counter(inst.JobName+"_errors_total", "Total job errors")

	wrapped := make(chan Job, cfg.Concurrency*2)
	go func() {
		for j := range jobsCh {
			wrapped <- &instrumentedJob{
				Job:        j,
				tracerName: inst.TracerName,
				latency:    latency.(prometheus.Observer),
				total:      total,
				errors:     errors,
			}
		}
		close(wrapped)
	}()

	return &InstrumentedPool{Pool: NewPool(cfg, wrapped)}
}

type instrumentedJob struct {
	Job        Job
	tracerName string
	latency    prometheus.Observer
	total      prometheus.Counter
	errors     prometheus.Counter
}

func (j *instrumentedJob) Run(ctx context.Context) error {
	start := time.Now()
	ctx, span := tracing.StartSpan(ctx, j.tracerName, "job")
	defer span.End()

	err := j.Job.Run(ctx)
	j.total.Inc()
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.SetAttributes(attribute.Bool("error", true))
		j.errors.Inc()
	}
	j.latency.Observe(time.Since(start).Seconds())
	return err
}
