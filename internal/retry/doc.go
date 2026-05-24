// Package retry implements a context-aware exponential back-off retry helper
// used by the webhook client to transparently re-attempt failed deliveries.
//
// Basic usage:
//
//	err := retry.Do(ctx, retry.DefaultConfig(), func() error {
//		return client.Send(payload)
//	})
//
// DefaultConfig provides four attempts with a 500 ms base delay capped at
// 30 seconds, which is appropriate for transient network failures without
// overwhelming the receiving endpoint.
package retry
