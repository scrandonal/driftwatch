package metrics_test

import (
	"testing"
	"time"

	"github.com/yourorg/driftwatch/internal/metrics"
)

func TestNew_InitialisesStartedAt(t *testing.T) {
	before := time.Now()
	c := metrics.New()
	after := time.Now()

	if c.StartedAt.Before(before) || c.StartedAt.After(after) {
		t.Errorf("StartedAt %v not between %v and %v", c.StartedAt, before, after)
	}
}

func TestCollector_Counters(t *testing.T) {
	c := metrics.New()

	c.FilesChecked.Add(5)
	c.ChangesFound.Add(2)
	c.AlertsSent.Add(2)
	c.AlertErrors.Add(1)
	c.WatchErrors.Add(3)

	snap := c.Snapshot()

	if snap.FilesChecked != 5 {
		t.Errorf("FilesChecked: want 5, got %d", snap.FilesChecked)
	}
	if snap.ChangesFound != 2 {
		t.Errorf("ChangesFound: want 2, got %d", snap.ChangesFound)
	}
	if snap.AlertsSent != 2 {
		t.Errorf("AlertsSent: want 2, got %d", snap.AlertsSent)
	}
	if snap.AlertErrors != 1 {
		t.Errorf("AlertErrors: want 1, got %d", snap.AlertErrors)
	}
	if snap.WatchErrors != 3 {
		t.Errorf("WatchErrors: want 3, got %d", snap.WatchErrors)
	}
}

func TestSnapshot_Uptime(t *testing.T) {
	c := metrics.New()
	time.Sleep(10 * time.Millisecond)
	snap := c.Snapshot()

	if snap.Uptime < 10*time.Millisecond {
		t.Errorf("Uptime too short: %v", snap.Uptime)
	}
}

func TestSnapshot_ZeroValues(t *testing.T) {
	c := metrics.New()
	snap := c.Snapshot()

	if snap.FilesChecked != 0 || snap.ChangesFound != 0 ||
		snap.AlertsSent != 0 || snap.AlertErrors != 0 || snap.WatchErrors != 0 {
		t.Error("expected all counters to be zero on a fresh collector")
	}
}
