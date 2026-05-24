package difflog_test

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/driftwatch/driftwatch/internal/difflog"
)

func TestNew_NilWriterReturnsError(t *testing.T) {
	_, err := difflog.New(nil)
	if err == nil {
		t.Fatal("expected error for nil writer, got nil")
	}
}

func TestNew_ValidWriter(t *testing.T) {
	var buf bytes.Buffer
	l, err := difflog.New(&buf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if l == nil {
		t.Fatal("expected non-nil logger")
	}
}

func TestRecord_WritesJSONLine(t *testing.T) {
	var buf bytes.Buffer
	l, _ := difflog.New(&buf)

	now := time.Now().UTC().Truncate(time.Second)
	e := difflog.Entry{
		Timestamp: now,
		Path:      "/etc/app/config.yaml",
		PrevHash:  "abc123",
		NewHash:   "def456",
	}
	if err := l.Record(e); err != nil {
		t.Fatalf("Record returned error: %v", err)
	}

	line := strings.TrimSpace(buf.String())
	var got difflog.Entry
	if err := json.Unmarshal([]byte(line), &got); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}
	if got.Path != e.Path {
		t.Errorf("path: got %q, want %q", got.Path, e.Path)
	}
	if got.PrevHash != e.PrevHash {
		t.Errorf("prev_hash: got %q, want %q", got.PrevHash, e.PrevHash)
	}
	if got.NewHash != e.NewHash {
		t.Errorf("new_hash: got %q, want %q", got.NewHash, e.NewHash)
	}
}

func TestRecord_ZeroTimestampFilledIn(t *testing.T) {
	var buf bytes.Buffer
	l, _ := difflog.New(&buf)

	before := time.Now().UTC().Add(-time.Second)
	if err := l.Record(difflog.Entry{Path: "/tmp/x", PrevHash: "a", NewHash: "b"}); err != nil {
		t.Fatalf("Record error: %v", err)
	}
	after := time.Now().UTC().Add(time.Second)

	var got difflog.Entry
	if err := json.Unmarshal(bytes.TrimSpace(buf.Bytes()), &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if got.Timestamp.Before(before) || got.Timestamp.After(after) {
		t.Errorf("timestamp %v not in expected range [%v, %v]", got.Timestamp, before, after)
	}
}

func TestRecord_MultipleEntries(t *testing.T) {
	var buf bytes.Buffer
	l, _ := difflog.New(&buf)

	for i := 0; i < 3; i++ {
		if err := l.Record(difflog.Entry{Path: "/f", PrevHash: "x", NewHash: "y"}); err != nil {
			t.Fatalf("Record[%d] error: %v", i, err)
		}
	}
	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(lines) != 3 {
		t.Errorf("expected 3 lines, got %d", len(lines))
	}
}
