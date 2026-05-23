// Package debounce provides a simple debouncer to coalesce rapid
// successive events into a single delayed notification.
package debounce

import (
	"sync"
	"time"
)

// Debouncer delays execution of a function until after a quiet period
// has elapsed since the last call.
type Debouncer struct {
	delay  time.Duration
	mu     sync.Mutex
	timer  *time.Timer
}

// New creates a new Debouncer with the given delay duration.
func New(delay time.Duration) *Debouncer {
	return &Debouncer{delay: delay}
}

// Call schedules fn to be called after the debounce delay.
// If Call is invoked again before the delay elapses, the timer resets.
func (d *Debouncer) Call(fn func()) {
	d.mu.Lock()
	defer d.mu.Unlock()

	if d.timer != nil {
		d.timer.Stop()
	}
	d.timer = time.AfterFunc(d.delay, func() {
		fn()
	})
}

// Stop cancels any pending debounced call.
func (d *Debouncer) Stop() {
	d.mu.Lock()
	defer d.mu.Unlock()

	if d.timer != nil {
		d.timer.Stop()
		d.timer = nil
	}
}
