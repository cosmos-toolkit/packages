// Package db provides connection bootstrap, pool and healthcheck.
// Driver (pgx, mysql, etc.) must be imported by the user.
package db

import (
	"context"
	"database/sql"
	"time"
)

// Config configures the connection.
type Config struct {
	Driver       string
	DSN          string
	MaxOpenConns int
	MaxIdleConns int
	MaxLifetime  time.Duration
}

// Open opens a connection to the database and configures the pool.
func Open(cfg Config) (*sql.DB, error) {
	db, err := sql.Open(cfg.Driver, cfg.DSN)
	if err != nil {
		return nil, err
	}
	if cfg.MaxOpenConns > 0 {
		db.SetMaxOpenConns(cfg.MaxOpenConns)
	}
	if cfg.MaxIdleConns > 0 {
		db.SetMaxIdleConns(cfg.MaxIdleConns)
	}
	if cfg.MaxLifetime > 0 {
		db.SetConnMaxLifetime(cfg.MaxLifetime)
	}
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, err
	}
	return db, nil
}

// Ping checks if the database is reachable (useful for healthcheck).
func Ping(ctx context.Context, db *sql.DB) error {
	if ctx == nil {
		return db.Ping()
	}
	return db.PingContext(ctx)
}
