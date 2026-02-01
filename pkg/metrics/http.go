// Package metrics exposes HTTP handler for Prometheus scrape.
package metrics

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Handler returns the HTTP handler for metrics (e.g. GET /metrics).
func Handler() http.Handler {
	return promhttp.Handler()
}
