// Package signalctx wraps the standard os/signal package to expose a
// context-based shutdown mechanism.
//
// # Usage
//
//		ctx, stop := signalctx.WithShutdown(context.Background())
//		defer stop()
//
//		// Pass ctx to long-running components; they will unblock when
//		// SIGINT or SIGTERM is received.
//
// The returned stop function must always be called (typically via defer) to
// release the underlying signal channel even if the context is already done.
package signalctx
