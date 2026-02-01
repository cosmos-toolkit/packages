// Package retry provides retry with exponential backoff, jitter and max attempts.
// Reusable outside queues (workers, HTTP client, etc.).
package retry

import (
	"context"
	"math/rand"
	"time"
)

// Config defines retry behaviour.
type Config struct {
	MaxAttempts int           // max attempts (>= 1)
	Initial     time.Duration // initial backoff
	MaxBackoff  time.Duration // backoff ceiling
	Multiplier  float64       // exponential factor (e.g. 2)
	Jitter      float64       // jitter 0..1 (e.g. 0.2 = ±20%)
}

// DefaultConfig returns a reasonable config (5 attempts, 100ms initial, 2x, 20% jitter).
func DefaultConfig() Config {
	return Config{
		MaxAttempts: 5,
		Initial:     100 * time.Millisecond,
		MaxBackoff:  30 * time.Second,
		Multiplier:  2,
		Jitter:      0.2,
	}
}

// Do runs fn until success, context cancellation or attempts exhausted.
// Returns the last error from fn.
func Do(ctx context.Context, cfg Config, fn func() error) error {
	var lastErr error
	backoff := cfg.Initial
	for attempt := 0; attempt < cfg.MaxAttempts; attempt++ {
		if err := ctx.Err(); err != nil {
			return err
		}
		lastErr = fn()
		if lastErr == nil {
			return nil
		}
		if attempt == cfg.MaxAttempts-1 {
			break
		}
		// backoff with jitter
		d := addJitter(backoff, cfg.Jitter)
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(d):
			// next backoff
		}
		if cfg.Multiplier > 0 {
			backoff = time.Duration(float64(backoff) * cfg.Multiplier)
			if cfg.MaxBackoff > 0 && backoff > cfg.MaxBackoff {
				backoff = cfg.MaxBackoff
			}
		}
	}
	return lastErr
}

func addJitter(d time.Duration, jitter float64) time.Duration {
	if jitter <= 0 || jitter > 1 {
		return d
	}
	half := float64(d) * jitter
	delta := (rand.Float64()*2 - 1) * half
	return d + time.Duration(delta)
}
