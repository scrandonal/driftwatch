package watcher

import (
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"sync"
	"time"
)

// FileState holds the last known state of a watched file.
type FileState struct {
	Path    string
	Hash    string
	ModTime time.Time
}

// ChangeEvent represents a detected change in a config file.
type ChangeEvent struct {
	Path    string
	OldHash string
	NewHash string
	DetectedAt time.Time
}

// Watcher monitors a set of file paths for changes.
type Watcher struct {
	mu       sync.Mutex
	paths    []string
	states   map[string]FileState
	Interval time.Duration
	Changes  chan ChangeEvent
	stopCh   chan struct{}
}

// New creates a new Watcher for the given file paths.
func New(paths []string, interval time.Duration) *Watcher {
	return &Watcher{
		paths:    paths,
		states:   make(map[string]FileState),
		Interval: interval,
		Changes:  make(chan ChangeEvent, 16),
		stopCh:   make(chan struct{}),
	}
}

// Start begins polling files for changes.
func (w *Watcher) Start() {
	go func() {
		ticker := time.NewTicker(w.Interval)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				w.poll()
			case <-w.stopCh:
				close(w.Changes)
				return
			}
		}
	}()
}

// Stop halts the watcher.
func (w *Watcher) Stop() {
	close(w.stopCh)
}

func (w *Watcher) poll() {
	w.mu.Lock()
	defer w.mu.Unlock()

	for _, path := range w.paths {
		hash, err := hashFile(path)
		if err != nil {
			continue
		}
		prev, seen := w.states[path]
		if !seen {
			w.states[path] = FileState{Path: path, Hash: hash, ModTime: time.Now()}
			continue
		}
		if prev.Hash != hash {
			w.Changes <- ChangeEvent{
				Path:       path,
				OldHash:    prev.Hash,
				NewHash:    hash,
				DetectedAt: time.Now(),
			}
			w.states[path] = FileState{Path: path, Hash: hash, ModTime: time.Now()}
		}
	}
}

func hashFile(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()
	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", h.Sum(nil)), nil
}
