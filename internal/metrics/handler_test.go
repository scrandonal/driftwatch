package metrics_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/yourorg/driftwatch/internal/metrics"
)

func TestHandler_ReturnsJSON(t *testing.T) {
	c := metrics.New()
	c.FilesChecked.Add(10)
	c.AlertsSent.Add(3)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/metrics", nil)

	metrics.Handler(c)(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("want 200, got %d", rec.Code)
	}

	var payload map[string]any
	if err := json.NewDecoder(rec.Body).Decode(&payload); err != nil {
		t.Fatalf("decode error: %v", err)
	}

	if got := payload["files_checked"]; got != float64(10) {
		t.Errorf("files_checked: want 10, got %v", got)
	}
	if got := payload["alerts_sent"]; got != float64(3) {
		t.Errorf("alerts_sent: want 3, got %v", got)
	}
}

func TestHandler_ContentTypeIsJSON(t *testing.T) {
	c := metrics.New()
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/metrics", nil)

	metrics.Handler(c)(rec, req)

	ct := rec.Header().Get("Content-Type")
	if ct != "application/json" {
		t.Errorf("Content-Type: want application/json, got %q", ct)
	}
}

func TestHandler_MethodNotAllowed(t *testing.T) {
	c := metrics.New()
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/metrics", nil)

	metrics.Handler(c)(rec, req)

	if rec.Code != http.StatusMethodNotAllowed {
		t.Errorf("want 405, got %d", rec.Code)
	}
}
