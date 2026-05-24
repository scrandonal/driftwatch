// Package notify provides a rate-limited, retrying notification pipeline
// that sits between the alerter and the webhook client.
package notify

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/example/driftwatch/internal/ratelimit"
	"github.com/example/driftwatch/internal/retry"
)

// Sender is anything that can deliver a notification payload.
type Sender interface {
	Send(ctx context.Context, payload []byte) error
}

// Pipeline wraps a Sender with rate-limiting and retry logic.
type Pipeline struct {
	sender  Sender
	limiter *ratelimit.Limiter
	retry   retry.Config
	log     *slog.Logger
}

// Config holds Pipeline construction options.
type Config struct {
	MaxPerWindow int
	Retry        retry.Config
}

// DefaultConfig returns sensible defaults.
func DefaultConfig() Config {
	return Config{
		MaxPerWindow: 10,
		Retry:        retry.DefaultConfig(),
	}
}

// New creates a Pipeline.
func New(sender Sender, cfg Config, log *slog.Logger) (*Pipeline, error) {
	if sender == nil {
		return nil, fmt.Errorf("notify: sender must not be nil")
	}
	if log == nil {
		log = slog.Default()
	}
	return &Pipeline{
		sender:  sender,
		limiter: ratelimit.New(cfg.MaxPerWindow),
		retry:   cfg.Retry,
		log:     log,
	}, nil
}

// Dispatch delivers payload through the rate-limiter and retry loop.
func (p *Pipeline) Dispatch(ctx context.Context, payload []byte) error {
	if !p.limiter.Allow() {
		remaining, _ := p.limiter.Remaining()
		p.log.Warn("notify: rate limit reached, dropping notification",
			"remaining", remaining)
		return fmt.Errorf("notify: rate limit exceeded")
	}

	return retry.Do(ctx, p.retry, func(ctx context.Context) error {
		err := p.sender.Send(ctx, payload)
		if err != nil {
			p.log.Warn("notify: send failed, will retry", "err", err)
		}
		return err
	})
}
