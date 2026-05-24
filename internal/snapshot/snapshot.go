// Package snapshot provides persistent storage of file hash snapshots,
// allowing driftwatch to detect changes across restarts.
package snapshot

import (
	"encoding/json"
	"os"
	"sync"
	"time"
)

// Entry records the last known hash and modification time for a watched file.
type Entry struct {
	Hash      string    `json:"hash"`
	ModifiedAt time.Time `json:"modified_at"`
}

// Store holds a map of file paths to their snapshot entries and persists
// them to a JSON file on disk.
type Store struct {
	mu       sync.RWMutex
	entries  map[string]Entry
	filePath string
}

// New loads an existing snapshot file from path, or returns an empty Store
// if the file does not yet exist.
func New(path string) (*Store, error) {
	s := &Store{
		entries:  make(map[string]Entry),
		filePath: path,
	}

	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return s, nil
	}
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(data, &s.entries); err != nil {
		return nil, err
	}
	return s, nil
}

// Get returns the Entry for the given file path and whether it was found.
func (s *Store) Get(filePath string) (Entry, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	e, ok := s.entries[filePath]
	return e, ok
}

// Set updates the Entry for the given file path and flushes to disk.
func (s *Store) Set(filePath string, entry Entry) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.entries[filePath] = entry
	return s.flush()
}

// flush writes the current entries to the snapshot file. Must be called
// with s.mu held for writing.
func (s *Store) flush() error {
	data, err := json.MarshalIndent(s.entries, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.filePath, data, 0o644)
}

// All returns a shallow copy of all stored entries.
func (s *Store) All() map[string]Entry {
	s.mu.RLock()
	defer s.mu.RUnlock()
	copy := make(map[string]Entry, len(s.entries))
	for k, v := range s.entries {
		copy[k] = v
	}
	return copy
}
