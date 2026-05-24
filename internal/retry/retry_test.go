package retry_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/yourusername/driftwatch/internal/retry"
)

var errTemp = errors.New("temporary error")

func TestDo_SucceedsOnFirstAttempt(t *testing.T) {
	calls := 0
	err := retry.Do(context.Background(), retry.DefaultConfig(), func() error {
		calls++
		return nil
	})
	if err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
	if calls != 1 {
		t.Fatalf("expected 1 call, got %d", calls)
	}
}

func TestDo_RetriesUntilSuccess(t *testing.T) {
	calls := 0
	cfg := retry.Config{MaxAttempts: 4, BaseDelay: time.Millisecond, MaxDelay: 10 * time.Millisecond}
	err := retry.Do(context.Background(), cfg, func() error {
		calls++
		if calls < 3 {
			return errTemp
		}
		return nil
	})
	if err != nil {
		t.Fatalf("expected nil after eventual success, got %v", err)
	}
	if calls != 3 {
		t.Fatalf("expected 3 calls, got %d", calls)
	}
}

func TestDo_ExhaustsAttempts(t *testing.T) {
	cfg := retry.Config{MaxAttempts: 3, BaseDelay: time.Millisecond, MaxDelay: 5 * time.Millisecond}
	calls := 0
	err := retry.Do(context.Background(), cfg, func() error {
		calls++
		return errTemp
	})
	if !errors.Is(err, retry.ErrMaxAttempts) {
		t.Fatalf("expected ErrMaxAttempts, got %v", err)
	}
	if !errors.Is(err, errTemp) {
		t.Fatalf("expected wrapped errTemp, got %v", err)
	}
	if calls != 3 {
		t.Fatalf("expected 3 calls, got %d", calls)
	}
}

func TestDo_ContextCancelled(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	calls := 0
	err := retry.Do(ctx, retry.DefaultConfig(), func() error {
		calls++
		return errTemp
	})
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("expected context.Canceled, got %v", err)
	}
	if calls != 0 {
		t.Fatalf("expected 0 calls, got %d", calls)
	}
}

func TestDo_ZeroMaxAttemptsTreatedAsOne(t *testing.T) {
	cfg := retry.Config{MaxAttempts: 0, BaseDelay: time.Millisecond, MaxDelay: time.Millisecond}
	calls := 0
	_ = retry.Do(context.Background(), cfg, func() error {
		calls++
		return errTemp
	})
	if calls != 1 {
		t.Fatalf("expected 1 call, got %d", calls)
	}
}
