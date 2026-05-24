// Package ratelimit provides a simple token-bucket style rate limiter
// used to throttle outgoing webhook notifications.
//
// A Limiter is created with a maximum number of events allowed per time
// window. Each call to Allow consumes one token; tokens are refilled
// automatically once the window elapses.
//
// Example usage:
//
//	limiter := ratelimit.New(5, time.Minute)
//	if limiter.Allow() {
//		// send alert
//	}
package ratelimit
