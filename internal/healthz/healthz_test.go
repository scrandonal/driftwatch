package healthz_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/yourusername/driftwatch/internal/healthz"
)

// stubProvider implements healthz.StatusProvider for testing.
type stubProvider struct {
	paths     int
	startedAt time.Time
}

func (s *stubProvider) WatchedPaths() int      { return s.paths }
func (s *stubProvider) StartedAt() time.Time   { return s.startedAt }

func newStub(paths int, ago time.Duration) *stubProvider {
	return &stubProvider{paths: paths, startedAt: time.Now().Add(-ago)}
}

func TestHandler_ReturnsOK(t *testing.T) {
	h := healthz.Handler(newStub(3, 5*time.Second))
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/healthz", nil))

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
}

func TestHandler_ContentTypeJSON(t *testing.T) {
	h := healthz.Handler(newStub(1, time.Second))
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/healthz", nil))

	ct := rec.Header().Get("Content-Type")
	if ct != "application/json" {
		t.Fatalf("unexpected Content-Type: %s", ct)
	}
}

func TestHandler_ResponseBody(t *testing.T) {
	stub := newStub(4, 10*time.Second)
	h := healthz.Handler(stub)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/healthz", nil))

	var resp healthz.Response
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if resp.Status != "ok" {
		t.Errorf("expected status ok, got %s", resp.Status)
	}
	if resp.WatchedPaths != 4 {
		t.Errorf("expected 4 watched paths, got %d", resp.WatchedPaths)
	}
	if resp.UptimeSeconds < 9 {
		t.Errorf("expected uptime >= 9s, got %f", resp.UptimeSeconds)
	}
}

func TestHandler_MethodNotAllowed(t *testing.T) {
	h := healthz.Handler(newStub(0, 0))
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodPost, "/healthz", nil))

	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected 405, got %d", rec.Code)
	}
}

func TestHandler_HeadAllowed(t *testing.T) {
	h := healthz.Handler(newStub(2, time.Second))
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodHead, "/healthz", nil))

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 for HEAD, got %d", rec.Code)
	}
}
