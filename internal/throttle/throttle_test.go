package throttle_test

import (
	"testing"
	"time"

	"github.com/user/driftwatch/internal/throttle"
)

func TestAllow_FirstCallAlwaysPermitted(t *testing.T) {
	th := throttle.New(100 * time.Millisecond)
	if !th.Allow("key") {
		t.Fatal("expected first call to be allowed")
	}
}

func TestAllow_SecondCallWithinIntervalDenied(t *testing.T) {
	th := throttle.New(100 * time.Millisecond)
	th.Allow("key")
	if th.Allow("key") {
		t.Fatal("expected second call within interval to be denied")
	}
}

func TestAllow_DifferentKeysAreIndependent(t *testing.T) {
	th := throttle.New(100 * time.Millisecond)
	th.Allow("a")
	if !th.Allow("b") {
		t.Fatal("expected different key to be allowed")
	}
}

func TestAllow_PermittedAfterInterval(t *testing.T) {
	now := time.Now()
	th := throttle.New(50 * time.Millisecond)

	// Inject a controllable clock.
	type clockable interface {
		SetNow(func() time.Time)
	}

	// Use real sleep to cross the interval boundary.
	th.Allow("key")
	time.Sleep(60 * time.Millisecond)
	_ = now
	if !th.Allow("key") {
		t.Fatal("expected call after interval to be allowed")
	}
}

func TestReset_AllowsImmediateRetry(t *testing.T) {
	th := throttle.New(10 * time.Second)
	th.Allow("key")
	th.Reset("key")
	if !th.Allow("key") {
		t.Fatal("expected allow after reset")
	}
}

func TestRemaining_ZeroWhenNotThrottled(t *testing.T) {
	th := throttle.New(100 * time.Millisecond)
	if r := th.Remaining("key"); r != 0 {
		t.Fatalf("expected 0 remaining for unknown key, got %v", r)
	}
}

func TestRemaining_PositiveWhenThrottled(t *testing.T) {
	th := throttle.New(10 * time.Second)
	th.Allow("key")
	if r := th.Remaining("key"); r <= 0 {
		t.Fatalf("expected positive remaining, got %v", r)
	}
}

func TestNew_PanicsOnZeroInterval(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic for zero interval")
		}
	}()
	throttle.New(0)
}
