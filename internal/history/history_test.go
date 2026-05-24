package history_test

import (
	"testing"
	"time"

	"github.com/yourusername/driftwatch/internal/history"
)

func TestNew_DefaultCapacity(t *testing.T) {
	r := history.New(0)
	if r.Len() != 0 {
		t.Fatalf("expected empty ring, got %d", r.Len())
	}
}

func TestRecord_StoresEntry(t *testing.T) {
	r := history.New(5)
	now := time.Now().UTC()
	r.Record("/etc/app.yaml", "abc123", now)

	if r.Len() != 1 {
		t.Fatalf("expected 1 entry, got %d", r.Len())
	}
	entries := r.Latest(1)
	if entries[0].Path != "/etc/app.yaml" {
		t.Errorf("unexpected path: %s", entries[0].Path)
	}
	if entries[0].ShortHash != "abc123" {
		t.Errorf("unexpected hash: %s", entries[0].ShortHash)
	}
}

func TestRecord_ZeroTimestampFilledIn(t *testing.T) {
	r := history.New(5)
	before := time.Now().UTC()
	r.Record("/etc/app.yaml", "abc123", time.Time{})
	after := time.Now().UTC()

	entries := r.Latest(1)
	if entries[0].DetectedAt.Before(before) || entries[0].DetectedAt.After(after) {
		t.Errorf("timestamp not filled in correctly: %v", entries[0].DetectedAt)
	}
}

func TestLatest_NewestFirst(t *testing.T) {
	r := history.New(10)
	paths := []string{"/a", "/b", "/c"}
	for _, p := range paths {
		r.Record(p, "hash", time.Now().UTC())
	}
	entries := r.Latest(3)
	if entries[0].Path != "/c" || entries[1].Path != "/b" || entries[2].Path != "/a" {
		t.Errorf("wrong order: %v", entries)
	}
}

func TestRing_OverwritesOldest(t *testing.T) {
	r := history.New(3)
	for i := 0; i < 5; i++ {
		r.Record("/file", "h", time.Now().UTC())
	}
	if r.Len() != 3 {
		t.Fatalf("expected cap 3, got %d", r.Len())
	}
}

func TestLatest_NLargerThanCount(t *testing.T) {
	r := history.New(10)
	r.Record("/x", "h1", time.Now().UTC())
	entries := r.Latest(100)
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
}

func TestLatest_EmptyRing(t *testing.T) {
	r := history.New(5)
	if r.Latest(3) != nil {
		t.Error("expected nil from empty ring")
	}
}
