package ratelimit_test

import (
	"testing"
	"time"

	"github.com/yourusername/driftwatch/internal/ratelimit"
)

func TestLimiter_AllowsUpToMax(t *testing.T) {
	l := ratelimit.New(3, time.Minute)

	for i := 0; i < 3; i++ {
		if !l.Allow() {
			t.Fatalf("expected Allow() == true on call %d", i+1)
		}
	}

	if l.Allow() {
		t.Fatal("expected Allow() == false after exhausting tokens")
	}
}

func TestLimiter_RefillsAfterWindow(t *testing.T) {
	l := ratelimit.New(2, 50*time.Millisecond)

	l.Allow()
	l.Allow()

	if l.Allow() {
		t.Fatal("expected rate limit to be exhausted")
	}

	time.Sleep(60 * time.Millisecond)

	if !l.Allow() {
		t.Fatal("expected tokens to be refilled after window")
	}
}

func TestLimiter_Remaining(t *testing.T) {
	l := ratelimit.New(5, time.Minute)

	if r := l.Remaining(); r != 5 {
		t.Fatalf("expected 5 remaining, got %d", r)
	}

	l.Allow()
	l.Allow()

	if r := l.Remaining(); r != 3 {
		t.Fatalf("expected 3 remaining, got %d", r)
	}
}

func TestLimiter_RemainingRefillsAfterWindow(t *testing.T) {
	l := ratelimit.New(2, 50*time.Millisecond)
	l.Allow()
	l.Allow()

	time.Sleep(60 * time.Millisecond)

	if r := l.Remaining(); r != 2 {
		t.Fatalf("expected 2 after window reset, got %d", r)
	}
}

func TestLimiter_ZeroMax(t *testing.T) {
	l := ratelimit.New(0, time.Minute)
	if l.Allow() {
		t.Fatal("expected Allow() == false when max is 0")
	}
}
