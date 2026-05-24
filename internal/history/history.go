// Package history maintains an in-memory ring buffer of recent change
// events so that the /metrics and /healthz handlers can expose a
// human-readable recent-activity feed without hitting persistent storage.
package history

import (
	"sync"
	"time"
)

// Entry represents a single recorded change event.
type Entry struct {
	Path      string    `json:"path"`
	ShortHash string    `json:"short_hash"`
	DetectedAt time.Time `json:"detected_at"`
}

// Ring is a fixed-capacity, thread-safe ring buffer of change entries.
type Ring struct {
	mu       sync.Mutex
	buf      []Entry
	cap      int
	writeIdx int
	count    int
}

// New returns a Ring with the given capacity. If cap is less than 1 it
// defaults to 10.
func New(cap int) *Ring {
	if cap < 1 {
		cap = 10
	}
	return &Ring{
		buf: make([]Entry, cap),
		cap: cap,
	}
}

// Record appends a new entry to the ring, overwriting the oldest when full.
func (r *Ring) Record(path, shortHash string, at time.Time) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if at.IsZero() {
		at = time.Now().UTC()
	}
	r.buf[r.writeIdx] = Entry{Path: path, ShortHash: shortHash, DetectedAt: at}
	r.writeIdx = (r.writeIdx + 1) % r.cap
	if r.count < r.cap {
		r.count++
	}
}

// Latest returns up to n most-recent entries, newest first.
// If n <= 0 all stored entries are returned.
func (r *Ring) Latest(n int) []Entry {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.count == 0 {
		return nil
	}
	if n <= 0 || n > r.count {
		n = r.count
	}
	out := make([]Entry, n)
	for i := 0; i < n; i++ {
		idx := (r.writeIdx - 1 - i + r.cap) % r.cap
		out[i] = r.buf[idx]
	}
	return out
}

// Len returns the number of entries currently stored.
func (r *Ring) Len() int {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.count
}
