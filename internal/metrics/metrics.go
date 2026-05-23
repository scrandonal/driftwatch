// Package metrics provides simple in-memory counters for tracking
// driftwatch runtime statistics such as files checked, alerts sent,
// and errors encountered.
package metrics

import (
	"sync"
	"sync/atomic"
	"time"
)

// Collector holds runtime counters for driftwatch operations.
type Collector struct {
	mu          sync.RWMutex
	StartedAt   time.Time

	FilesChecked  atomic.Int64
	ChangesFound  atomic.Int64
	AlertsSent    atomic.Int64
	AlertErrors   atomic.Int64
	WatchErrors   atomic.Int64
}

// New creates and returns a new Collector initialised with the current time.
func New() *Collector {
	return &Collector{
		StartedAt: time.Now(),
	}
}

// Snapshot is a point-in-time copy of the collector's counters.
type Snapshot struct {
	Uptime        time.Duration
	FilesChecked  int64
	ChangesFound  int64
	AlertsSent    int64
	AlertErrors   int64
	WatchErrors   int64
}

// Snapshot returns a consistent read of all counters.
func (c *Collector) Snapshot() Snapshot {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return Snapshot{
		Uptime:       time.Since(c.StartedAt),
		FilesChecked: c.FilesChecked.Load(),
		ChangesFound: c.ChangesFound.Load(),
		AlertsSent:   c.AlertsSent.Load(),
		AlertErrors:  c.AlertErrors.Load(),
		WatchErrors:  c.WatchErrors.Load(),
	}
}
