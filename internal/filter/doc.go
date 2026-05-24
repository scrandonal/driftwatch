// Package filter provides glob-based path filtering for the driftwatch
// file watcher. It allows callers to specify include and exclude glob
// patterns so that only relevant file-system events are forwarded for
// further processing.
//
// A Filter is constructed once with [New] and is safe for concurrent use
// after construction. Patterns follow the syntax accepted by
// [path/filepath.Match].
package filter
