package filter_test

import (
	"testing"

	"github.com/yourusername/driftwatch/internal/filter"
)

func TestFilter_Allow_NoGlobs_AcceptsAll(t *testing.T) {
	f := filter.New(nil, nil)
	paths := []string{"/etc/app.yaml", "/tmp/foo.log", "config/db.toml"}
	for _, p := range paths {
		if !f.Allow(p) {
			t.Errorf("expected Allow(%q) = true, got false", p)
		}
	}
}

func TestFilter_Allow_IncludeGlob_FiltersCorrectly(t *testing.T) {
	f := filter.New([]string{"*.yaml", "*.yml"}, nil)

	allowed := []string{"config.yaml", "app.yml", "/etc/service.yaml"}
	for _, p := range allowed {
		if !f.Allow(p) {
			t.Errorf("expected Allow(%q) = true", p)
		}
	}

	denied := []string{"config.toml", "/etc/hosts", "notes.txt"}
	for _, p := range denied {
		if f.Allow(p) {
			t.Errorf("expected Allow(%q) = false", p)
		}
	}
}

func TestFilter_Allow_ExcludeGlob_TakesPrecedence(t *testing.T) {
	// Include all yaml but exclude anything under /tmp.
	f := filter.New([]string{"*.yaml"}, []string{"/tmp/*"})

	if !f.Allow("config.yaml") {
		t.Error("expected config.yaml to be allowed")
	}
	if f.Allow("/tmp/config.yaml") {
		t.Error("expected /tmp/config.yaml to be denied by exclude glob")
	}
}

func TestFilter_Allow_ExcludeOnly_DeniesMatches(t *testing.T) {
	f := filter.New(nil, []string{"*.log", "*.tmp"})

	if f.Allow("app.log") {
		t.Error("expected app.log to be denied")
	}
	if f.Allow("/var/run/session.tmp") {
		t.Error("expected session.tmp to be denied")
	}
	if !f.Allow("/etc/app.yaml") {
		t.Error("expected app.yaml to be allowed")
	}
}

func TestFilter_Allow_BasenameMatching(t *testing.T) {
	// Pattern without path separator should match on basename.
	f := filter.New([]string{"*.yaml"}, nil)

	if !f.Allow("/deeply/nested/path/config.yaml") {
		t.Error("expected basename match to allow /deeply/nested/path/config.yaml")
	}
}

func TestFilter_Allow_InvalidGlob_Skipped(t *testing.T) {
	// An invalid glob pattern (bracket not closed) must not panic.
	f := filter.New([]string{"[invalid"}, nil)
	// With only an invalid include glob the file should not be allowed.
	if f.Allow("config.yaml") {
		t.Error("expected no match when only include glob is invalid")
	}
}
