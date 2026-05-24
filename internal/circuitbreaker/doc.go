// Package circuitbreaker provides a thread-safe circuit breaker implementation
// for use in driftwatch's webhook notification pipeline.
//
// The circuit breaker transitions between three states:
//
//   - Closed: normal operation; calls are allowed through.
//   - Open: the downstream service is considered unhealthy; calls are
//     rejected immediately with ErrOpen to avoid cascading failures.
//   - Half-Open: after the reset timeout elapses, one probe call is
//     permitted. A success closes the circuit; a failure re-opens it.
//
// Example usage:
//
//	cb := circuitbreaker.New(circuitbreaker.DefaultConfig())
//
//	if err := cb.Allow(); err != nil {
//		// circuit is open — skip the call
//		return err
//	}
//	if err := doRequest(); err != nil {
//		cb.RecordFailure()
//		return err
//	}
//	cb.RecordSuccess()
package circuitbreaker
