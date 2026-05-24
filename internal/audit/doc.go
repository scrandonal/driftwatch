// Package audit provides structured append-only audit logging for driftwatch.
//
// Each time a monitored file changes, an audit record is written as a newline-
// delimited JSON entry containing the file path, old and new content hashes,
// a short hash for human-readable display, and the UTC timestamp of the event.
//
// # Writers
//
// An audit.Logger can be constructed around any io.Writer, making it easy to
// direct output to stdout, a rotating file, or a test buffer:
//
//	logger, err := audit.New(os.Stdout)
//
// For production use, NewFile opens (or creates) a file in append mode so that
// records survive process restarts without overwriting previous history:
//
//	logger, err := audit.NewFile("/var/log/driftwatch/audit.jsonl")
//
// # Record format
//
// Each line written by Logger.Record is a self-contained JSON object:
//
//	{"ts":"2024-01-15T10:04:05Z","path":"/etc/app/config.yaml",
//	 "old_hash":"abc12345","new_hash":"def67890"}
//
// A zero-value timestamp in the supplied Record is automatically replaced with
// the current UTC time before serialisation.
package audit
