package version_test

import (
	"strings"
	"testing"

	"github.com/yourusername/driftwatch/internal/version"
)

func TestGet_ReturnsDefaults(t *testing.T) {
	info := version.Get()

	if info.Version == "" {
		t.Error("expected Version to be non-empty")
	}
	if info.Commit == "" {
		t.Error("expected Commit to be non-empty")
	}
	if info.Date == "" {
		t.Error("expected Date to be non-empty")
	}
}

func TestGet_DefaultValues(t *testing.T) {
	info := version.Get()

	// In test builds the ldflags are not set, so defaults apply.
	if info.Version != "dev" {
		t.Errorf("expected default Version \"dev\", got %q", info.Version)
	}
	if info.Commit != "unknown" {
		t.Errorf("expected default Commit \"unknown\", got %q", info.Commit)
	}
	if info.Date != "unknown" {
		t.Errorf("expected default Date \"unknown\", got %q", info.Date)
	}
}

func TestInfo_String_ContainsVersion(t *testing.T) {
	info := version.Get()
	s := info.String()

	if !strings.Contains(s, "driftwatch") {
		t.Errorf("String() missing \"driftwatch\": %q", s)
	}
	if !strings.Contains(s, info.Version) {
		t.Errorf("String() missing version %q: %q", info.Version, s)
	}
	if !strings.Contains(s, info.Commit) {
		t.Errorf("String() missing commit %q: %q", info.Commit, s)
	}
	if !strings.Contains(s, info.Date) {
		t.Errorf("String() missing date %q: %q", info.Date, s)
	}
}

func TestInfo_String_Format(t *testing.T) {
	info := version.Info{
		Version: "1.2.3",
		Commit:  "abc1234",
		Date:    "2024-01-15T10:00:00Z",
	}

	got := info.String()
	want := "driftwatch 1.2.3 (commit=abc1234, built=2024-01-15T10:00:00Z)"

	if got != want {
		t.Errorf("String() = %q, want %q", got, want)
	}
}
