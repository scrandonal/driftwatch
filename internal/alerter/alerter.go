// Package alerter coordinates detection events and webhook delivery.
package alerter

import (
	"context"
	"fmt"
	"log"
	"time"
)

// Event represents a detected configuration file change.
type Event struct {
	FilePath  string    `json:"file_path"`
	OldHash   string    `json:"old_hash"`
	NewHash   string    `json:"new_hash"`
	DetectedAt time.Time `json:"detected_at"`
}

// Sender is the interface for delivering alert payloads.
type Sender interface {
	Send(ctx context.Context, payload any) error
}

// Alerter dispatches change events to a webhook sender.
type Alerter struct {
	sender Sender
	logger *log.Logger
}

// New creates a new Alerter with the given Sender and logger.
func New(sender Sender, logger *log.Logger) *Alerter {
	if logger == nil {
		logger = log.Default()
	}
	return &Alerter{sender: sender, logger: logger}
}

// Notify sends an alert for the provided change Event.
// It returns an error if delivery fails.
func (a *Alerter) Notify(ctx context.Context, event Event) error {
	a.logger.Printf("alerter: change detected in %s (old=%s new=%s)",
		event.FilePath, shortHash(event.OldHash), shortHash(event.NewHash))

	if err := a.sender.Send(ctx, event); err != nil {
		return fmt.Errorf("alerter: failed to send notification for %s: %w", event.FilePath, err)
	}

	a.logger.Printf("alerter: notification delivered for %s", event.FilePath)
	return nil
}

// shortHash returns the first 8 characters of a hash for log readability.
func shortHash(h string) string {
	if len(h) <= 8 {
		return h
	}
	return h[:8]
}
