// Package ratelimit provides a simple token-bucket rate limiter
// used to suppress alert floods when many files change at once.
package ratelimit

import (
	"sync"
	"time"
)

// Limiter is a thread-safe token-bucket rate limiter.
type Limiter struct {
	mu       sync.Mutex
	tokens   int
	max      int
	refillAt time.Time
	window   time.Duration
}

// New creates a Limiter that allows at most max events per window duration.
// Tokens are fully replenished after each window elapses.
func New(max int, window time.Duration) *Limiter {
	return &Limiter{
		tokens:   max,
		max:      max,
		window:   window,
		refillAt: time.Now().Add(window),
	}
}

// Allow reports whether an event should be allowed through.
// It consumes one token; if none remain before the window resets,
// the call returns false.
func (l *Limiter) Allow() bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := time.Now()
	if now.After(l.refillAt) {
		l.tokens = l.max
		l.refillAt = now.Add(l.window)
	}

	if l.tokens <= 0 {
		return false
	}
	l.tokens--
	return true
}

// Remaining returns the number of tokens left in the current window.
func (l *Limiter) Remaining() int {
	l.mu.Lock()
	defer l.mu.Unlock()

	if time.Now().After(l.refillAt) {
		return l.max
	}
	return l.tokens
}
