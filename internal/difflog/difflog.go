// Package difflog records file diff events to a structured append-only log.
// Each entry captures the file path, old and new hashes, and a timestamp so
// that operators can reconstruct a full change history for any watched file.
package difflog

import (
	"encoding/json"
	"fmt"
	"io"
	"time"
)

// Entry represents a single recorded change event.
type Entry struct {
	// Timestamp is when the change was detected (UTC).
	Timestamp time.Time `json:"timestamp"`
	// Path is the absolute path of the file that changed.
	Path string `json:"path"`
	// OldHash is the SHA-256 hex digest before the change, or empty string on
	// first observation.
	OldHash string `json:"old_hash"`
	// NewHash is the SHA-256 hex digest after the change.
	NewHash string `json:"new_hash"`
}

// Logger writes diff entries as newline-delimited JSON to an underlying writer.
type Logger struct {
	w   io.Writer
	enc *json.Encoder
}

// New creates a Logger that writes to w. It returns an error if w is nil.
func New(w io.Writer) (*Logger, error) {
	if w == nil {
		return nil, fmt.Errorf("difflog: writer must not be nil")
	}
	enc := json.NewEncoder(w)
	enc.SetEscapeHTML(false)
	return &Logger{w: w, enc: enc}, nil
}

// Record appends a change entry to the log. If e.Timestamp is zero it is
// replaced with the current UTC time before writing.
func (l *Logger) Record(e Entry) error {
	if e.Timestamp.IsZero() {
		e.Timestamp = time.Now().UTC()
	}
	if err := l.enc.Encode(e); err != nil {
		return fmt.Errorf("difflog: encode entry: %w", err)
	}
	return nil
}
