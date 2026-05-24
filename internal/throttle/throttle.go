// Package throttle provides a token-bucket style throttle that limits
// how frequently a keyed action may be performed within a sliding window.
package throttle

import (
	"sync"
	"time"
)

// Throttle tracks per-key timestamps and enforces a minimum interval
// between successive calls for the same key.
type Throttle struct {
	mu       sync.Mutex
	last     map[string]time.Time
	interval time.Duration
	now      func() time.Time
}

// New returns a Throttle that allows at most one action per key within
// the given interval. interval must be positive; if it is not, New panics.
func New(interval time.Duration) *Throttle {
	if interval <= 0 {
		panic("throttle: interval must be positive")
	}
	return &Throttle{
		last:     make(map[string]time.Time),
		interval: interval,
		now:      time.Now,
	}
}

// Allow reports whether the action identified by key is permitted at the
// current time. If it is permitted the internal timestamp is updated so
// that subsequent calls within the interval are denied.
func (t *Throttle) Allow(key string) bool {
	t.mu.Lock()
	defer t.mu.Unlock()

	now := t.now()
	if last, ok := t.last[key]; ok && now.Sub(last) < t.interval {
		return false
	}
	t.last[key] = now
	return true
}

// Reset removes the recorded timestamp for key, allowing the next call to
// Allow for that key to succeed immediately.
func (t *Throttle) Reset(key string) {
	t.mu.Lock()
	defer t.mu.Unlock()
	delete(t.last, key)
}

// Remaining returns the duration until key is next allowed. It returns 0
// if the key is not currently throttled.
func (t *Throttle) Remaining(key string) time.Duration {
	t.mu.Lock()
	defer t.mu.Unlock()

	last, ok := t.last[key]
	if !ok {
		return 0
	}
	remaining := t.interval - t.now().Sub(last)
	if remaining < 0 {
		return 0
	}
	return remaining
}
