package circuitbreaker_test

import (
	"testing"
	"time"

	"github.com/yourusername/driftwatch/internal/circuitbreaker"
)

func newBreaker(maxFailures int, resetTimeout time.Duration) *circuitbreaker.Breaker {
	return circuitbreaker.New(circuitbreaker.Config{
		MaxFailures:  maxFailures,
		ResetTimeout: resetTimeout,
	})
}

func TestBreaker_AllowsWhenClosed(t *testing.T) {
	b := newBreaker(3, time.Second)
	if err := b.Allow(); err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
}

func TestBreaker_OpensAfterMaxFailures(t *testing.T) {
	b := newBreaker(3, time.Second)
	for i := 0; i < 3; i++ {
		b.RecordFailure()
	}
	if b.State() != circuitbreaker.StateOpen {
		t.Fatal("expected state to be Open")
	}
	if err := b.Allow(); err != circuitbreaker.ErrOpen {
		t.Fatalf("expected ErrOpen, got %v", err)
	}
}

func TestBreaker_SuccessResetsBreakerFromHalfOpen(t *testing.T) {
	b := newBreaker(2, 10*time.Millisecond)
	b.RecordFailure()
	b.RecordFailure()

	time.Sleep(20 * time.Millisecond)

	// Should transition to half-open
	if err := b.Allow(); err != nil {
		t.Fatalf("expected nil in half-open, got %v", err)
	}
	b.RecordSuccess()
	if b.State() != circuitbreaker.StateClosed {
		t.Fatal("expected state to be Closed after success")
	}
}

func TestBreaker_RemainsOpenBeforeTimeout(t *testing.T) {
	b := newBreaker(1, 10*time.Second)
	b.RecordFailure()
	if err := b.Allow(); err != circuitbreaker.ErrOpen {
		t.Fatalf("expected ErrOpen before timeout, got %v", err)
	}
}

func TestBreaker_DefaultConfigApplied(t *testing.T) {
	b := circuitbreaker.New(circuitbreaker.Config{})
	if b == nil {
		t.Fatal("expected non-nil breaker")
	}
	if b.State() != circuitbreaker.StateClosed {
		t.Fatal("expected initial state to be Closed")
	}
}

func TestBreaker_PartialFailuresDontOpen(t *testing.T) {
	b := newBreaker(5, time.Second)
	for i := 0; i < 4; i++ {
		b.RecordFailure()
	}
	if b.State() != circuitbreaker.StateClosed {
		t.Fatal("expected state to remain Closed below threshold")
	}
	if err := b.Allow(); err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
}
