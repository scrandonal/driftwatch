package notify_test

import (
	"context"
	"errors"
	"log/slog"
	"testing"
	"time"

	"github.com/example/driftwatch/internal/notify"
	"github.com/example/driftwatch/internal/retry"
)

type mockSender struct {
	calls  int
	failN  int // fail the first N calls
	lastPayload []byte
}

func (m *mockSender) Send(_ context.Context, payload []byte) error {
	m.calls++
	m.lastPayload = payload
	if m.calls <= m.failN {
		return errors.New("transient error")
	}
	return nil
}

func fastRetry() retry.Config {
	return retry.Config{
		MaxAttempts: 3,
		InitialDelay: time.Millisecond,
		MaxDelay: 5 * time.Millisecond,
		Multiplier: 1.5,
	}
}

func TestPipeline_Dispatch_Success(t *testing.T) {
	sender := &mockSender{}
	cfg := notify.DefaultConfig()
	cfg.Retry = fastRetry()
	p, err := notify.New(sender, cfg, slog.Default())
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	if err := p.Dispatch(context.Background(), []byte(`{"file":"a.yaml"}`)); err != nil {
		t.Fatalf("Dispatch: %v", err)
	}
	if sender.calls != 1 {
		t.Errorf("expected 1 call, got %d", sender.calls)
	}
}

func TestPipeline_Dispatch_RetriesOnTransient(t *testing.T) {
	sender := &mockSender{failN: 2}
	cfg := notify.DefaultConfig()
	cfg.Retry = fastRetry()
	p, _ := notify.New(sender, cfg, slog.Default())

	if err := p.Dispatch(context.Background(), []byte(`{}`)); err != nil {
		t.Fatalf("expected eventual success, got: %v", err)
	}
	if sender.calls != 3 {
		t.Errorf("expected 3 calls (2 failures + 1 success), got %d", sender.calls)
	}
}

func TestPipeline_Dispatch_RateLimitExceeded(t *testing.T) {
	sender := &mockSender{}
	cfg := notify.Config{
		MaxPerWindow: 2,
		Retry:        fastRetry(),
	}
	p, _ := notify.New(sender, cfg, slog.Default())

	for i := 0; i < 2; i++ {
		if err := p.Dispatch(context.Background(), []byte(`{}`)); err != nil {
			t.Fatalf("call %d: unexpected error: %v", i, err)
		}
	}
	if err := p.Dispatch(context.Background(), []byte(`{}`)); err == nil {
		t.Fatal("expected rate limit error on third call")
	}
}

func TestNew_NilSenderReturnsError(t *testing.T) {
	_, err := notify.New(nil, notify.DefaultConfig(), nil)
	if err == nil {
		t.Fatal("expected error for nil sender")
	}
}
