package snapshot_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/yourusername/driftwatch/internal/snapshot"
)

func tempPath(t *testing.T) string {
	t.Helper()
	return filepath.Join(t.TempDir(), "snapshot.json")
}

func TestNew_EmptyWhenFileAbsent(t *testing.T) {
	s, err := snapshot.New(tempPath(t))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := s.All(); len(got) != 0 {
		t.Fatalf("expected empty store, got %d entries", len(got))
	}
}

func TestSet_PersistsAndGet(t *testing.T) {
	path := tempPath(t)
	s, _ := snapshot.New(path)

	entry := snapshot.Entry{Hash: "abc123", ModifiedAt: time.Now().Truncate(time.Second)}
	if err := s.Set("/etc/app.yaml", entry); err != nil {
		t.Fatalf("Set failed: %v", err)
	}

	got, ok := s.Get("/etc/app.yaml")
	if !ok {
		t.Fatal("expected entry to exist")
	}
	if got.Hash != entry.Hash {
		t.Errorf("hash mismatch: got %q want %q", got.Hash, entry.Hash)
	}
}

func TestNew_LoadsExistingFile(t *testing.T) {
	path := tempPath(t)

	// Write initial data via a store.
	s1, _ := snapshot.New(path)
	_ = s1.Set("/etc/app.yaml", snapshot.Entry{Hash: "deadbeef", ModifiedAt: time.Now()})

	// Re-open from same path.
	s2, err := snapshot.New(path)
	if err != nil {
		t.Fatalf("reload failed: %v", err)
	}
	e, ok := s2.Get("/etc/app.yaml")
	if !ok {
		t.Fatal("entry missing after reload")
	}
	if e.Hash != "deadbeef" {
		t.Errorf("unexpected hash after reload: %q", e.Hash)
	}
}

func TestNew_ErrorOnCorruptFile(t *testing.T) {
	path := tempPath(t)
	if err := os.WriteFile(path, []byte("not json{"), 0o644); err != nil {
		t.Fatalf("setup: %v", err)
	}
	_, err := snapshot.New(path)
	if err == nil {
		t.Fatal("expected error for corrupt snapshot file")
	}
}

func TestAll_ReturnsCopy(t *testing.T) {
	s, _ := snapshot.New(tempPath(t))
	_ = s.Set("/a", snapshot.Entry{Hash: "h1"})
	_ = s.Set("/b", snapshot.Entry{Hash: "h2"})

	all := s.All()
	if len(all) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(all))
	}
	// Mutating the returned map must not affect the store.
	delete(all, "/a")
	if _, ok := s.Get("/a"); !ok {
		t.Error("store was mutated by modifying All() result")
	}
}
