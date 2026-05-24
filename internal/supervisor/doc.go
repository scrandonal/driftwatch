// Package supervisor provides a simple worker supervisor that runs a
// worker function in a loop, restarting it automatically if it returns
// a non-nil error. The supervisor respects context cancellation: when
// the supplied context is cancelled the current worker run is allowed
// to finish and the supervisor exits without restarting.
//
// Usage:
//
//	sv := supervisor.New(func(ctx context.Context) error {
//		// long-running work here
//		return nil
//	})
//	sv.Run(ctx)
//
// The supervisor applies a small back-off delay between restarts to
// avoid tight error loops that would hammer downstream resources.
package supervisor
