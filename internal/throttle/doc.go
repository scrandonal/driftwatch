// Package throttle implements a simple per-key throttle that enforces a
// minimum interval between successive actions for the same key.
//
// Unlike ratelimit, which uses a token-bucket model with a maximum burst
// count, throttle is designed for scenarios where exactly one event per
// key should be processed within a given window — for example, suppressing
// duplicate webhook deliveries for the same config file path.
//
// Basic usage:
//
//	th := throttle.New(30 * time.Second)
//
//	if th.Allow(filePath) {
//		// send alert
//	} else {
//		// skip — already alerted recently
//	}
//
// Reset clears the recorded timestamp for a key, allowing the next call
// to Allow to succeed immediately regardless of the interval.
package throttle
