// Package circuitbreaker implements a simple circuit breaker pattern
// that prevents repeated calls to a failing downstream service.
package circuitbreaker

import (
	"errors"
	"sync"
	"time"
)

// ErrOpen is returned when the circuit breaker is in the open state.
var ErrOpen = errors.New("circuit breaker is open")

// State represents the current state of the circuit breaker.
type State int

const (
	StateClosed State = iota
	StateOpen
	StateHalfOpen
)

// Config holds configuration for the circuit breaker.
type Config struct {
	// MaxFailures is the number of consecutive failures before opening.
	MaxFailures int
	// ResetTimeout is how long to wait before transitioning to half-open.
	ResetTimeout time.Duration
}

// DefaultConfig returns a Config with sensible defaults.
func DefaultConfig() Config {
	return Config{
		MaxFailures:  5,
		ResetTimeout: 30 * time.Second,
	}
}

// Breaker is a circuit breaker.
type Breaker struct {
	mu          sync.Mutex
	cfg         Config
	state       State
	failures    int
	lastFailure time.Time
}

// New creates a new Breaker with the given config.
func New(cfg Config) *Breaker {
	if cfg.MaxFailures <= 0 {
		cfg.MaxFailures = DefaultConfig().MaxFailures
	}
	if cfg.ResetTimeout <= 0 {
		cfg.ResetTimeout = DefaultConfig().ResetTimeout
	}
	return &Breaker{cfg: cfg}
}

// Allow reports whether a call should be allowed through.
// It returns ErrOpen when the circuit is open.
func (b *Breaker) Allow() error {
	b.mu.Lock()
	defer b.mu.Unlock()

	switch b.state {
	case StateOpen:
		if time.Since(b.lastFailure) >= b.cfg.ResetTimeout {
			b.state = StateHalfOpen
			return nil
		}
		return ErrOpen
	default:
		return nil
	}
}

// RecordSuccess records a successful call, resetting the breaker.
func (b *Breaker) RecordSuccess() {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.failures = 0
	b.state = StateClosed
}

// RecordFailure records a failed call, potentially opening the circuit.
func (b *Breaker) RecordFailure() {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.failures++
	b.lastFailure = time.Now()
	if b.failures >= b.cfg.MaxFailures {
		b.state = StateOpen
	}
}

// State returns the current state of the breaker.
func (b *Breaker) State() State {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.state
}
