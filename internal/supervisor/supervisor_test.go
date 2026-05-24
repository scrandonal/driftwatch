package supervisor_test

import (
	"context"
	"sync/atomic"
	"testing"
	"time"

	"github.com/yourusername/driftwatch/internal/supervisor"
)

func TestSupervisor_StartsAndStops(t *testing.T) {
	var callCount atomic.Int32

	sv := supervisor.New(func(ctx context.Context) error {
		callCount.Add(1)
		<-ctx.Done()
		return nil
	})

	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	done := make(chan struct{})
	go func() {
		sv.Run(ctx)
		close(done)
	}()

	select {
	case <-done:
		// ok
	case <-time.After(500 * time.Millisecond):
		t.Fatal("supervisor did not stop within timeout")
	}

	if callCount.Load() == 0 {
		t.Error("expected worker to be called at least once")
	}
}

func TestSupervisor_RestartsOnError(t *testing.T) {
	var callCount atomic.Int32

	sv := supervisor.New(func(ctx context.Context) error {
		n := callCount.Add(1)
		if n < 3 {
			return fmt.Errorf("simulated failure %d", n)
		}
		<-ctx.Done()
		return nil
	})

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	done := make(chan struct{})
	go func() {
		sv.Run(ctx)
		close(done)
	}()

	// Wait until worker has been called at least 3 times
	deadline := time.Now().Add(1 * time.Second)
	for time.Now().Before(deadline) {
		if callCount.Load() >= 3 {
			break
		}
		time.Sleep(20 * time.Millisecond)
	}

	if callCount.Load() < 3 {
		t.Errorf("expected at least 3 restarts, got %d", callCount.Load())
	}

	cancel()
	select {
	case <-done:
	case <-time.After(500 * time.Millisecond):
		t.Fatal("supervisor did not stop after cancel")
	}
}

func TestSupervisor_ContextCancelledImmediately(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // already cancelled

	sv := supervisor.New(func(ctx context.Context) error {
		<-ctx.Done()
		return nil
	})

	done := make(chan struct{})
	go func() {
		sv.Run(ctx)
		close(done)
	}()

	select {
	case <-done:
		// ok
	case <-time.After(300 * time.Millisecond):
		t.Fatal("supervisor should exit quickly when context is pre-cancelled")
	}
}
