package dispatcher_test

import (
	"testing"
	"time"

	"github.com/example/driftwatch/internal/dispatcher"
)

func TestDefaultBuildConfig_Defaults(t *testing.T) {
	bc := dispatcher.DefaultBuildConfig()
	if bc.DebounceDuration != 200*time.Millisecond {
		t.Errorf("expected 200ms debounce, got %v", bc.DebounceDuration)
	}
	if len(bc.IncludeGlobs) != 0 {
		t.Errorf("expected no include globs by default")
	}
	if len(bc.ExcludeGlobs) != 0 {
		t.Errorf("expected no exclude globs by default")
	}
}

func TestBuild_NoGlobsNilFilter(t *testing.T) {
	sender := &mockSender{}
	bc := dispatcher.BuildConfig{DebounceDuration: 0}
	cfg := dispatcher.Build(bc, sender)
	if cfg.Filter != nil {
		t.Error("expected nil filter when no globs provided")
	}
	if cfg.Debounce != nil {
		t.Error("expected nil debouncer when duration is zero")
	}
}

func TestBuild_WithGlobsCreatesFilter(t *testing.T) {
	sender := &mockSender{}
	bc := dispatcher.BuildConfig{
		IncludeGlobs:     []string{"*.yaml"},
		DebounceDuration: 0,
	}
	cfg := dispatcher.Build(bc, sender)
	if cfg.Filter == nil {
		t.Fatal("expected non-nil filter when globs provided")
	}
}

func TestBuild_WithDebounceDurationCreatesDebouncer(t *testing.T) {
	sender := &mockSender{}
	bc := dispatcher.BuildConfig{
		DebounceDuration: 100 * time.Millisecond,
	}
	cfg := dispatcher.Build(bc, sender)
	if cfg.Debounce == nil {
		t.Fatal("expected non-nil debouncer when duration > 0")
	}
}
