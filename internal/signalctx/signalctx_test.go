package signalctx

import (
	"context"
	"os"
	"syscall"
	"testing"
	"time"
)

func TestWithShutdown_CancelledBySignal(t *testing.T) {
	t.Parallel()

	ctx, stop := withSignals(context.Background(), syscall.SIGUSR1)
	defer stop()

	// Send the signal to ourselves.
	p, err := os.FindProcess(os.Getpid())
	if err != nil {
		t.Fatalf("find process: %v", err)
	}
	if err := p.Signal(syscall.SIGUSR1); err != nil {
		t.Fatalf("send signal: %v", err)
	}

	select {
	case <-ctx.Done():
		// expected
	case <-time.After(2 * time.Second):
		t.Fatal("context was not cancelled after signal")
	}
}

func TestWithShutdown_StopCancelsContext(t *testing.T) {
	t.Parallel()

	ctx, stop := withSignals(context.Background(), syscall.SIGUSR1)

	stop()

	select {
	case <-ctx.Done():
		// expected
	case <-time.After(time.Second):
		t.Fatal("context was not cancelled after stop()")
	}
}

func TestWithShutdown_ParentCancelPropagates(t *testing.T) {
	t.Parallel()

	parent, cancelParent := context.WithCancel(context.Background())
	ctx, stop := withSignals(parent, syscall.SIGUSR1)
	defer stop()

	cancelParent()

	select {
	case <-ctx.Done():
		// expected
	case <-time.After(time.Second):
		t.Fatal("child context not cancelled when parent was cancelled")
	}
}

func TestWithShutdown_NotCancelledWithoutSignal(t *testing.T) {
	t.Parallel()

	ctx, stop := withSignals(context.Background(), syscall.SIGUSR1)
	defer stop()

	select {
	case <-ctx.Done():
		t.Fatal("context cancelled unexpectedly")
	case <-time.After(50 * time.Millisecond):
		// expected — no signal sent, context should remain open
	}
}
