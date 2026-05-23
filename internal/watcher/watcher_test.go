package watcher

import (
	"os"
	"testing"
	"time"
)

func writeTempFile(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp("", "driftwatch-*.conf")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}
	f.Close()
	return f.Name()
}

func TestWatcher_DetectsChange(t *testing.T) {
	path := writeTempFile(t, "initial content")
	defer os.Remove(path)

	w := New([]string{path}, 50*time.Millisecond)
	w.Start()
	defer w.Stop()

	// Allow one poll to record the initial state.
	time.Sleep(100 * time.Millisecond)

	// Modify the file.
	if err := os.WriteFile(path, []byte("changed content"), 0644); err != nil {
		t.Fatalf("failed to update file: %v", err)
	}

	select {
	case event := <-w.Changes:
		if event.Path != path {
			t.Errorf("expected path %s, got %s", path, event.Path)
		}
		if event.OldHash == event.NewHash {
			t.Error("expected hashes to differ after file change")
		}
	case <-time.After(500 * time.Millisecond):
		t.Fatal("timed out waiting for change event")
	}
}

func TestWatcher_NoChangeNoEvent(t *testing.T) {
	path := writeTempFile(t, "stable content")
	defer os.Remove(path)

	w := New([]string{path}, 50*time.Millisecond)
	w.Start()
	defer w.Stop()

	time.Sleep(200 * time.Millisecond)

	select {
	case event := <-w.Changes:
		t.Errorf("unexpected change event: %+v", event)
	default:
		// pass — no events expected
	}
}

func TestHashFile_ConsistentHash(t *testing.T) {
	path := writeTempFile(t, "consistent")
	defer os.Remove(path)

	h1, err := hashFile(path)
	if err != nil {
		t.Fatalf("hashFile error: %v", err)
	}
	h2, err := hashFile(path)
	if err != nil {
		t.Fatalf("hashFile error: %v", err)
	}
	if h1 != h2 {
		t.Errorf("expected identical hashes, got %s and %s", h1, h2)
	}
}
