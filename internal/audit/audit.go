// Package audit provides a simple append-only audit log that records
// every file-change event processed by driftwatch.
package audit

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sync"
	"time"
)

// Entry represents a single audit-log record.
type Entry struct {
	Timestamp time.Time `json:"timestamp"`
	Path      string    `json:"path"`
	OldHash   string    `json:"old_hash"`
	NewHash   string    `json:"new_hash"`
	Alerted   bool      `json:"alerted"`
}

// Logger writes audit entries to an io.Writer (typically a file).
type Logger struct {
	mu  sync.Mutex
	out io.Writer
}

// New returns a Logger that writes to w.
// Pass os.Stdout or an *os.File opened with os.O_APPEND.
func New(w io.Writer) (*Logger, error) {
	if w == nil {
		return nil, fmt.Errorf("audit: writer must not be nil")
	}
	return &Logger{out: w}, nil
}

// NewFile opens (or creates) the file at path for appending and returns a
// Logger backed by that file together with a close function.
func NewFile(path string) (*Logger, func() error, error) {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		return nil, nil, fmt.Errorf("audit: open %s: %w", path, err)
	}
	l, err := New(f)
	if err != nil {
		_ = f.Close()
		return nil, nil, err
	}
	return l, f.Close, nil
}

// Record serialises e as a JSON line and appends it to the underlying writer.
// It is safe for concurrent use.
func (l *Logger) Record(e Entry) error {
	if e.Timestamp.IsZero() {
		e.Timestamp = time.Now().UTC()
	}
	b, err := json.Marshal(e)
	if err != nil {
		return fmt.Errorf("audit: marshal entry: %w", err)
	}
	b = append(b, '\n')

	l.mu.Lock()
	defer l.mu.Unlock()
	_, err = l.out.Write(b)
	return err
}
