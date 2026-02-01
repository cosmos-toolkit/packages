// Package metrics provides Prometheus helpers: latency, error rate, throughput.
package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// Counter creates a counter registered in the default registry (promauto).
func Counter(name, help string) prometheus.Counter {
	return promauto.NewCounter(prometheus.CounterOpts{
		Name: name,
		Help: help,
	})
}

// Histogram creates a histogram with default buckets (latency in seconds).
func Histogram(name, help string, buckets []float64) prometheus.Histogram {
	if len(buckets) == 0 {
		buckets = []float64{.005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5, 10}
	}
	return promauto.NewHistogram(prometheus.HistogramOpts{
		Name:    name,
		Help:    help,
		Buckets: buckets,
	})
}

// DefaultLatencyBuckets returns typical buckets for HTTP latency (seconds).
func DefaultLatencyBuckets() []float64 {
	return []float64{.005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5, 10}
}
