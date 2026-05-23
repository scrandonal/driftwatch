// Package debounce provides a Debouncer that coalesces multiple rapid
// successive triggers into a single delayed function call.
//
// It is used by driftwatch to prevent alert storms when a config file
// is written in multiple small chunks (e.g. by editors that perform
// atomic saves via rename). Instead of firing a webhook for every
// intermediate filesystem event, the Debouncer waits for a configurable
// quiet period before invoking the alerter.
//
// Basic usage:
//
//	d := debounce.New(500 * time.Millisecond)
//	// Inside the watcher event loop:
//	d.Call(func() { alerter.Notify(path, newHash) })
package debounce
