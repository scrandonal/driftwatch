// Package ratelimit implements a simple token-bucket rate limiter for
// controlling the frequency of outbound webhook alerts in driftwatch.
//
// A Limiter is created with a maximum number of allowed events and a
// refill window. Each call to Allow consumes one token; once tokens are
// exhausted the limiter returns false until the window elapses and the
// bucket is replenished to its maximum capacity.
//
// Example usage:
//
//	limiter := ratelimit.New(10, time.Minute)
//	if limiter.Allow() {
//		// send alert
//	}
package ratelimit
