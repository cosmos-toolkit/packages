// Package logger fornece um wrapper sobre slog com suporte a contexto
// (trace_id, request_id) e modos CLI / API / Worker.
package logger

import (
	"context"
	"log/slog"
	"os"

	"github.com/cosmos-toolkit/pkgs/pkg/contextx"
)

// Mode define o formato de saída do logger.
type Mode string

const (
	ModeCLI    Mode = "cli"    // humano, texto
	ModeAPI    Mode = "api"    // JSON, para HTTP
	ModeWorker Mode = "worker" // JSON, para jobs
)

// Logger encapsula slog com campos de contexto e modo.
type Logger struct {
	*slog.Logger
	mode Mode
}

// Config configura o logger.
type Config struct {
	Mode  Mode
	Level slog.Level
}

// DefaultConfig retorna configuração padrão (API, Info).
func DefaultConfig() Config {
	return Config{Mode: ModeAPI, Level: slog.LevelInfo}
}

// New cria um logger com a configuração dada.
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

// WithContext retorna um logger que inclui trace_id, request_id e user_id
// do contexto (via contextx), quando presentes.
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

// Mode retorna o modo do logger.
func (l *Logger) Mode() Mode { return l.mode }
