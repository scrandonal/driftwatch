// Package retry provides a simple exponential back-off retry mechanism
// for transient failures such as webhook delivery errors.
package retry

import (
	"context"
	"errors"
	"math"
	"time"
)

// ErrMaxAttempts is returned when all retry attempts are exhausted.
var ErrMaxAttempts = errors.New("retry: max attempts reached")

// Config holds the parameters that control retry behaviour.
type Config struct {
	// MaxAttempts is the total number of attempts (including the first).
	MaxAttempts int
	// BaseDelay is the delay before the second attempt.
	BaseDelay time.Duration
	// MaxDelay caps the computed exponential delay.
	MaxDelay time.Duration
}

// DefaultConfig returns a Config suitable for webhook delivery.
func DefaultConfig() Config {
	return Config{
		MaxAttempts: 4,
		BaseDelay:   500 * time.Millisecond,
		MaxDelay:    30 * time.Second,
	}
}

// Do calls fn up to cfg.MaxAttempts times, backing off exponentially
// between attempts. It stops early if ctx is cancelled or fn returns nil.
func Do(ctx context.Context, cfg Config, fn func() error) error {
	if cfg.MaxAttempts <= 0 {
		cfg.MaxAttempts = 1
	}

	var last error
	for attempt := 0; attempt < cfg.MaxAttempts; attempt++ {
		if err := ctx.Err(); err != nil {
			return err
		}

		last = fn()
		if last == nil {
			return nil
		}

		if attempt == cfg.MaxAttempts-1 {
			break
		}

		delay := delay(cfg.BaseDelay, cfg.MaxDelay, attempt)
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(delay):
		}
	}

	return errors.Join(ErrMaxAttempts, last)
}

// delay computes the capped exponential back-off for a given attempt index.
func delay(base, max time.Duration, attempt int) time.Duration {
	d := time.Duration(float64(base) * math.Pow(2, float64(attempt)))
	if d > max {
		return max
	}
	return d
}
