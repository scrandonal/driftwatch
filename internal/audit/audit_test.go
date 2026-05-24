package audit_test

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/yourusername/driftwatch/internal/audit"
)

func TestNew_NilWriterReturnsError(t *testing.T) {
	_, err := audit.New(nil)
	if err == nil {
		t.Fatal("expected error for nil writer, got nil")
	}
}

func TestRecord_WritesJSONLine(t *testing.T) {
	var buf bytes.Buffer
	l, err := audit.New(&buf)
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	now := time.Now().UTC().Truncate(time.Second)
	e := audit.Entry{
		Timestamp: now,
		Path:      "/etc/app/config.yaml",
		OldHash:   "aabbcc",
		NewHash:   "ddeeff",
		Alerted:   true,
	}
	if err := l.Record(e); err != nil {
		t.Fatalf("Record: %v", err)
	}

	line := strings.TrimSpace(buf.String())
	var got audit.Entry
	if err := json.Unmarshal([]byte(line), &got); err != nil {
		t.Fatalf("unmarshal output: %v", err)
	}
	if got.Path != e.Path {
		t.Errorf("path: want %q, got %q", e.Path, got.Path)
	}
	if got.OldHash != e.OldHash || got.NewHash != e.NewHash {
		t.Errorf("hashes mismatch")
	}
	if !got.Alerted {
		t.Error("alerted: want true, got false")
	}
}

func TestRecord_ZeroTimestampFilledIn(t *testing.T) {
	var buf bytes.Buffer
	l, _ := audit.New(&buf)

	before := time.Now().UTC()
	_ = l.Record(audit.Entry{Path: "/tmp/x"})
	after := time.Now().UTC()

	var got audit.Entry
	_ = json.Unmarshal(bytes.TrimSpace(buf.Bytes()), &got)
	if got.Timestamp.Before(before) || got.Timestamp.After(after) {
		t.Errorf("timestamp %v not in expected range [%v, %v]", got.Timestamp, before, after)
	}
}

func TestNewFile_CreatesAndAppends(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "audit.log")

	l, close, err := audit.NewFile(path)
	if err != nil {
		t.Fatalf("NewFile: %v", err)
	}
	defer close()

	_ = l.Record(audit.Entry{Path: "/a", NewHash: "123"})
	_ = l.Record(audit.Entry{Path: "/b", NewHash: "456"})
	_ = close()

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile: %v", err)
	}
	lines := strings.Split(strings.TrimSpace(string(data)), "\n")
	if len(lines) != 2 {
		t.Errorf("want 2 lines, got %d", len(lines))
	}
}
