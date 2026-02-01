// Package httpx provides server bootstrap, graceful shutdown and default healthcheck.
package httpx

import (
	"context"
	"net/http"
	"time"
)

// Server wraps http.Server with graceful shutdown.
type Server struct {
	*http.Server
}

// NewServer creates an HTTP server with default timeouts.
func NewServer(addr string, handler http.Handler) *Server {
	return &Server{
		Server: &http.Server{
			Addr:         addr,
			Handler:      handler,
			ReadTimeout:  15 * time.Second,
			WriteTimeout: 15 * time.Second,
			IdleTimeout:  60 * time.Second,
		},
	}
}

// Run starts the server and blocks until ctx is cancelled; then performs graceful shutdown.
func (s *Server) Run(ctx context.Context) error {
	done := make(chan error, 1)
	go func() { done <- s.ListenAndServe() }()

	select {
	case err := <-done:
		return err
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := s.Shutdown(shutdownCtx); err != nil {
			return err
		}
		return ctx.Err()
	}
}

// HealthHandler returns a handler that responds 200 OK at /health (or the given path).
func HealthHandler(path string) http.HandlerFunc {
	if path == "" {
		path = "/health"
	}
	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != path {
			http.NotFound(w, r)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}
}
