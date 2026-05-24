// Package supervisor wires together the watcher, debouncer, alerter, and
// rate-limiter into a single long-running component.
package supervisor

import (
	"context"
	"log/slog"
	"time"

	"github.com/yourusername/driftwatch/internal/alerter"
	"github.com/yourusername/driftwatch/internal/debounce"
	"github.com/yourusername/driftwatch/internal/metrics"
	"github.com/yourusername/driftwatch/internal/ratelimit"
	"github.com/yourusername/driftwatch/internal/watcher"
)

// Supervisor coordinates file watching and alert dispatching.
type Supervisor struct {
	watcher   *watcher.Watcher
	alerter   *alerter.Alerter
	limiter   *ratelimit.Limiter
	debouncer *debounce.Debouncer
	collector *metrics.Collector
	logger    *slog.Logger
}

// Config holds the parameters needed to build a Supervisor.
type Config struct {
	Watcher       *watcher.Watcher
	Alerter       *alerter.Alerter
	Limiter       *ratelimit.Limiter
	DebounceDelay time.Duration
	Collector     *metrics.Collector
	Logger        *slog.Logger
}

// New creates a Supervisor from the provided Config.
func New(cfg Config) *Supervisor {
	s := &Supervisor{
		watcher:   cfg.Watcher,
		alerter:   cfg.Alerter,
		limiter:   cfg.Limiter,
		collector: cfg.Collector,
		logger:    cfg.Logger,
	}
	s.debouncer = debounce.New(cfg.DebounceDelay, s.flush)
	return s
}

// Run starts the supervision loop and blocks until ctx is cancelled.
func (s *Supervisor) Run(ctx context.Context) error {
	events := s.watcher.Events()
	for {
		select {
		case <-ctx.Done():
			s.debouncer.Stop()
			return ctx.Err()
		case ev, ok := <-events:
			if !ok {
				return nil
			}
			s.logger.Info("change detected", "path", ev.Path)
			s.debouncer.Call()
		}
	}
}

// flush is invoked by the debouncer after the quiet period expires.
func (s *Supervisor) flush() {
	if !s.limiter.Allow() {
		s.logger.Warn("rate limit reached; skipping alert",
			"remaining", s.limiter.Remaining())
		s.collector.IncDropped()
		return
	}
	if err := s.alerter.Notify(context.Background()); err != nil {
		s.logger.Error("alert delivery failed", "err", err)
		s.collector.IncErrors()
		return
	}
	s.collector.IncAlerts()
}
