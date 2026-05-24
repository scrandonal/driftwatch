// Package dispatcher coordinates file-change events from the watcher
// and fans them out to the notify pipeline, applying filtering and
// debouncing before forwarding.
package dispatcher

import (
	"context"
	"log/slog"

	"github.com/example/driftwatch/internal/debounce"
	"github.com/example/driftwatch/internal/filter"
	"github.com/example/driftwatch/internal/notify"
)

// Event represents a detected file change.
type Event struct {
	Path    string
	OldHash string
	NewHash string
}

// Sender abstracts the notify pipeline for testing.
type Sender interface {
	Dispatch(ctx context.Context, path, oldHash, newHash string) error
}

// Dispatcher fans change events to the notify pipeline.
type Dispatcher struct {
	filter   *filter.Filter
	debounce *debounce.Debouncer
	sender   Sender
	log      *slog.Logger
	events   chan Event
}

// Config holds Dispatcher configuration.
type Config struct {
	Filter   *filter.Filter
	Debounce *debounce.Debouncer
	Sender   Sender
	Log      *slog.Logger
}

// New creates a Dispatcher ready to receive events.
func New(cfg Config) *Dispatcher {
	if cfg.Log == nil {
		cfg.Log = slog.Default()
	}
	return &Dispatcher{
		filter:   cfg.Filter,
		debounce: cfg.Debounce,
		sender:   cfg.Sender,
		log:      cfg.Log,
		events:   make(chan Event, 64),
	}
}

// Send enqueues a change event for processing.
func (d *Dispatcher) Send(e Event) {
	select {
	case d.events <- e:
	default:
		d.log.Warn("dispatcher: event queue full, dropping event", "path", e.Path)
	}
}

// Run processes events until ctx is cancelled.
func (d *Dispatcher) Run(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case e := <-d.events:
			if d.filter != nil && !d.filter.Allow(e.Path) {
				d.log.Debug("dispatcher: event filtered", "path", e.Path)
				continue
			}
			d.forward(ctx, e)
		}
	}
}

func (d *Dispatcher) forward(ctx context.Context, e Event) {
	if err := d.sender.Dispatch(ctx, e.Path, e.OldHash, e.NewHash); err != nil {
		d.log.Error("dispatcher: failed to dispatch event", "path", e.Path, "err", err)
	}
}

// Ensure notify.Pipeline satisfies Sender at compile time.
var _ Sender = (*notify.Pipeline)(nil)
