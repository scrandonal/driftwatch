package debounce_test

import (
	"sync/atomic"
	"testing"
	"time"

	"github.com/yourorg/driftwatch/internal/debounce"
)

func TestDebouncer_CallsOnceAfterQuietPeriod(t *testing.T) {
	var count int32
	d := debounce.New(50 * time.Millisecond)

	// Fire multiple times in rapid succession.
	for i := 0; i < 5; i++ {
		d.Call(func() {
			atomic.AddInt32(&count, 1)
		})
		time.Sleep(10 * time.Millisecond)
	}

	// Wait for the debounce delay to expire.
	time.Sleep(100 * time.Millisecond)

	if got := atomic.LoadInt32(&count); got != 1 {
		t.Errorf("expected fn to be called once, got %d", got)
	}
}

func TestDebouncer_Stop_CancelsPendingCall(t *testing.T) {
	var count int32
	d := debounce.New(80 * time.Millisecond)

	d.Call(func() {
		atomic.AddInt32(&count, 1)
	})

	// Stop before the delay elapses.
	d.Stop()
	time.Sleep(120 * time.Millisecond)

	if got := atomic.LoadInt32(&count); got != 0 {
		t.Errorf("expected fn not to be called after Stop, got %d calls", got)
	}
}

func TestDebouncer_MultipleCallsResetTimer(t *testing.T) {
	var timestamps []time.Time
	var mu sync.Mutex

	d := debounce.New(60 * time.Millisecond)

	start := time.Now()
	for i := 0; i < 3; i++ {
		d.Call(func() {
			mu.Lock()
			timestamps = append(timestamps, time.Now())
			mu.Unlock()
		})
		time.Sleep(20 * time.Millisecond)
	}

	time.Sleep(120 * time.Millisecond)

	mu.Lock()
	defer mu.Unlock()

	if len(timestamps) != 1 {
		t.Fatalf("expected 1 call, got %d", len(timestamps))
	}

	// The call should happen at least 60 ms after the last Call invocation.
	elapsed := timestamps[0].Sub(start)
	if elapsed < 100*time.Millisecond {
		t.Errorf("expected debounced call to be delayed, elapsed=%v", elapsed)
	}
}
