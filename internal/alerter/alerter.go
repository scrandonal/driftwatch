// Package alerter composes the webhook client with optional rate-limiting
// and structured logging to dispatch drift notifications.
package alerter

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/yourusername/driftwatch/internal/ratelimit"
	"github.com/yourusername/driftwatch/internal/webhook"
)

// Sender is the interface used to dispatch alert payloads.
type Sender interface {
	Send(ctx context.Context, payload any) error
}

// Alerter sends drift notifications through a Sender, gated by a rate limiter.
type Alerter struct {
	sender  Sender
	limiter *ratelimit.Limiter
	logger  *slog.Logger
}

// New constructs an Alerter. logger must not be nil.
func New(sender Sender, limiter *ratelimit.Limiter, logger *slog.Logger) *Alerter {
	if logger == nil {
		panic("alerter: logger must not be nil")
	}
	return &Alerter{sender: sender, limiter: limiter, logger: logger}
}

// payload is the JSON body sent to the webhook endpoint.
type payload struct {
	Event    string `json:"event"`
	Path     string `json:"path"`
	Hash     string `json:"hash"`
	Detected string `json:"detected_at"`
}

// Notify dispatches a drift alert for the given file path and its new hash.
// If the rate limiter is set and exhausted, the notification is dropped and
// a warning is logged instead.
func (a *Alerter) Notify(ctx context.Context, path, hash string) error {
	if a.limiter != nil && !a.limiter.Allow() {
		a.logger.Warn("alert suppressed by rate limiter",
			"path", path,
			"remaining", a.limiter.Remaining(),
		)
		return nil
	}

	p := payload{
		Event:    "drift.detected",
		Path:     path,
		Hash:     shortHash(hash),
		Detected: time.Now().UTC().Format(time.RFC3339),
	}

	if err := a.sender.Send(ctx, p); err != nil {
		return fmt.Errorf("alerter: send failed for %q: %w", path, err)
	}

	a.logger.Info("drift alert sent", "path", path, "hash", p.Hash)
	return nil
}

func shortHash(h string) string {
	if len(h) > 12 {
		return h[:12]
	}
	return h
}

// compile-time check
var _ Sender = (*webhook.Client)(nil)
