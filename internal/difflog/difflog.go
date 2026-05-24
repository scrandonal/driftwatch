// Package difflog records a human-readable log of file change events,
// capturing the file path, previous hash, new hash, and timestamp.
package difflog

import (
	"encoding/json"
	"fmt"
	"io"
	"time"
)

// Entry represents a single file-change record written to the diff log.
type Entry struct {
	Timestamp time.Time `json:"timestamp"`
	Path      string    `json:"path"`
	PrevHash  string    `json:"prev_hash"`
	NewHash   string    `json:"new_hash"`
}

// Logger writes diff entries as newline-delimited JSON to an io.Writer.
type Logger struct {
	w io.Writer
}

// New returns a Logger that writes to w.
// It returns an error if w is nil.
func New(w io.Writer) (*Logger, error) {
	if w == nil {
		return nil, fmt.Errorf("difflog: writer must not be nil")
	}
	return &Logger{w: w}, nil
}

// Record encodes e as a JSON line and writes it to the underlying writer.
// If e.Timestamp is zero it is set to the current UTC time.
func (l *Logger) Record(e Entry) error {
	if e.Timestamp.IsZero() {
		e.Timestamp = time.Now().UTC()
	}
	b, err := json.Marshal(e)
	if err != nil {
		return fmt.Errorf("difflog: marshal entry: %w", err)
	}
	b = append(b, '\n')
	if _, err := l.w.Write(b); err != nil {
		return fmt.Errorf("difflog: write entry: %w", err)
	}
	return nil
}
