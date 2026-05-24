package alerter_test

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"testing"
	"time"

	"github.com/yourusername/driftwatch/internal/alerter"
	"github.com/yourusername/driftwatch/internal/ratelimit"
)

type mockSender struct {
	called  int
	lastCtx context.Context
	lastPay any
	err     error
}

func (m *mockSender) Send(ctx context.Context, payload any) error {
	m.called++
	m.lastCtx = ctx
	m.lastPay = payload
	return m.err
}

func newTestAlerter(sender alerter.Sender, limiter *ratelimit.Limiter) *alerter.Alerter {
	logger := slog.New(slog.NewTextHandler(os.Stderr, nil))
	return alerter.New(sender, limiter, logger)
}

func TestAlerter_Notify_Success(t *testing.T) {
	sender := &mockSender{}
	a := newTestAlerter(sender, nil)

	if err := a.Notify(context.Background(), "/etc/app.conf", "abc123"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if sender.called != 1 {
		t.Fatalf("expected sender called once, got %d", sender.called)
	}
}

func TestAlerter_Notify_SenderError(t *testing.T) {
	sender := &mockSender{err: errors.New("connection refused")}
	a := newTestAlerter(sender, nil)

	err := a.Notify(context.Background(), "/etc/app.conf", "abc123")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestAlerter_New_NilLogger(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic for nil logger")
		}
	}()
	alerter.New(&mockSender{}, nil, nil)
}

func TestAlerter_Notify_RateLimited(t *testing.T) {
	sender := &mockSender{}
	limiter := ratelimit.New(2, time.Minute)
	a := newTestAlerter(sender, limiter)

	for i := 0; i < 5; i++ {
		if err := a.Notify(context.Background(), "/etc/app.conf", "abc"); err != nil {
			t.Fatalf("unexpected error on call %d: %v", i, err)
		}
	}

	if sender.called != 2 {
		t.Fatalf("expected 2 sends (rate limit=2), got %d", sender.called)
	}
}
