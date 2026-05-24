package healthz

import (
	"sync/atomic"
	"time"
)

// DefaultProvider is a thread-safe implementation of StatusProvider that
// driftwatch's main package can use directly.
type DefaultProvider struct {
	startedAt    time.Time
	watchedPaths atomic.Int64
}

// NewDefaultProvider creates a DefaultProvider whose start time is set to now.
func NewDefaultProvider() *DefaultProvider {
	return &DefaultProvider{startedAt: time.Now()}
}

// StartedAt returns the time the provider (and therefore the service) was
// created.
func (p *DefaultProvider) StartedAt() time.Time {
	return p.startedAt
}

// WatchedPaths returns the current number of watched paths.
func (p *DefaultProvider) WatchedPaths() int {
	return int(p.watchedPaths.Load())
}

// SetWatchedPaths updates the count of paths currently under observation.
// It is safe to call from multiple goroutines.
func (p *DefaultProvider) SetWatchedPaths(n int) {
	p.watchedPaths.Store(int64(n))
}
