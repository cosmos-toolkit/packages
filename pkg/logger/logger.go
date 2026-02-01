// Package logger provides a slog wrapper with context support
// (trace_id, request_id) and CLI / API / Worker modes.
package logger

import (
	"context"
	"log/slog"
	"os"

	"github.com/cosmos-toolkit/pkgs/pkg/contextx"
)

// Mode defines the logger output format.
type Mode string

const (
	ModeCLI    Mode = "cli"    // human-readable text
	ModeAPI    Mode = "api"    // JSON for HTTP
	ModeWorker Mode = "worker" // JSON for jobs
)

// Logger wraps slog with context fields and mode.
type Logger struct {
	*slog.Logger
	mode Mode
}

// Config configures the logger.
type Config struct {
	Mode  Mode
	Level slog.Level
}

// DefaultConfig returns default configuration (API, Info).
func DefaultConfig() Config {
	return Config{Mode: ModeAPI, Level: slog.LevelInfo}
}

// New creates a logger with the given configuration.
func New(cfg Config) *Logger {
	var handler slog.Handler
	opts := &slog.HandlerOptions{Level: cfg.Level}
	if cfg.Mode == ModeCLI {
		handler = slog.NewTextHandler(os.Stdout, opts)
	} else {
		handler = slog.NewJSONHandler(os.Stdout, opts)
	}
	return &Logger{Logger: slog.New(handler), mode: cfg.Mode}
}

// WithContext returns a logger that includes trace_id, request_id and user_id
// from the context (via contextx) when present.
func (l *Logger) WithContext(ctx context.Context) *slog.Logger {
	if ctx == nil {
		return l.Logger
	}
	logger := l.Logger
	if s := contextx.TraceID(ctx); s != "" {
		logger = logger.With(slog.String("trace_id", s))
	}
	if s := contextx.RequestID(ctx); s != "" {
		logger = logger.With(slog.String("request_id", s))
	}
	if s := contextx.UserID(ctx); s != "" {
		logger = logger.With(slog.String("user_id", s))
	}
	return logger
}

// Mode returns the logger mode.
func (l *Logger) Mode() Mode { return l.mode }
