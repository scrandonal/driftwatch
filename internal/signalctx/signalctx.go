// Package signalctx provides a context that is cancelled when the process
// receives an OS interrupt or termination signal.
package signalctx

import (
	"context"
	"os"
	"os/signal"
	"syscall"
)

// Signals is the default set of OS signals that trigger cancellation.
var Signals = []os.Signal{
	syscall.SIGINT,
	syscall.SIGTERM,
}

// WithShutdown returns a context that is cancelled when one of the watched
// signals is received, along with a stop function that releases resources.
//
// The caller must invoke stop when the context is no longer needed to avoid
// a goroutine leak, even if the context has already been cancelled.
//
//		ctx, stop := signalctx.WithShutdown(context.Background())
//		defer stop()
func WithShutdown(parent context.Context) (context.Context, context.CancelFunc) {
	return withSignals(parent, Signals...)
}

// withSignals is the testable core — it accepts an explicit signal list.
func withSignals(parent context.Context, sigs ...os.Signal) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(parent)

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, sigs...)

	go func() {
		select {
		case <-ch:
			cancel()
		case <-ctx.Done():
		}
		signal.Stop(ch)
	}()

	stop := func() {
		cancel()
		signal.Stop(ch)
	}

	return ctx, stop
}
