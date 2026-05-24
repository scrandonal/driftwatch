// Package history provides a thread-safe, fixed-capacity ring buffer that
// records recent file-change events detected by driftwatch.
//
// # Overview
//
// A [Ring] holds up to a configurable number of [Entry] values. When the
// buffer is full the oldest entry is silently overwritten, keeping memory
// usage bounded regardless of how long the process runs.
//
// # HTTP handler
//
// [Handler] wraps a Ring and exposes its contents as a JSON array over HTTP.
// Callers may pass a `?limit=N` query parameter (capped at 100) to control
// how many entries are returned. Entries are ordered newest-first.
//
// # Usage
//
//	r := history.New(50)
//	r.Record("/etc/app.yaml", "a1b2c3d", time.Now().UTC())
//	entries := r.Latest(10)
//
// Typical integration: call Record from the alerter or dispatcher whenever a
// change event is forwarded to the webhook pipeline.
package history
