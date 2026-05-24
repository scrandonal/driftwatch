package history_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/yourusername/driftwatch/internal/history"
)

func newRingWithEntries(n int) *history.Ring {
	r := history.New(50)
	for i := 0; i < n; i++ {
		r.Record("/file", "hash", time.Now().UTC())
	}
	return r
}

func TestHandler_ReturnsJSON(t *testing.T) {
	r := newRingWithEntries(3)
	h := history.Handler(r)

	req := httptest.NewRequest(http.MethodGet, "/history", nil)
	rw := httptest.NewRecorder()
	h(rw, req)

	if rw.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rw.Code)
	}
	var entries []history.Entry
	if err := json.NewDecoder(rw.Body).Decode(&entries); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if len(entries) != 3 {
		t.Errorf("expected 3 entries, got %d", len(entries))
	}
}

func TestHandler_ContentTypeIsJSON(t *testing.T) {
	h := history.Handler(history.New(5))
	req := httptest.NewRequest(http.MethodGet, "/history", nil)
	rw := httptest.NewRecorder()
	h(rw, req)

	ct := rw.Header().Get("Content-Type")
	if ct != "application/json" {
		t.Errorf("unexpected content-type: %s", ct)
	}
}

func TestHandler_MethodNotAllowed(t *testing.T) {
	h := history.Handler(history.New(5))
	req := httptest.NewRequest(http.MethodPost, "/history", nil)
	rw := httptest.NewRecorder()
	h(rw, req)

	if rw.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected 405, got %d", rw.Code)
	}
}

func TestHandler_LimitQueryParam(t *testing.T) {
	r := newRingWithEntries(10)
	h := history.Handler(r)

	req := httptest.NewRequest(http.MethodGet, "/history?limit=3", nil)
	rw := httptest.NewRecorder()
	h(rw, req)

	var entries []history.Entry
	_ = json.NewDecoder(rw.Body).Decode(&entries)
	if len(entries) != 3 {
		t.Errorf("expected 3 entries with limit=3, got %d", len(entries))
	}
}

func TestHandler_EmptyRingReturnsEmptyArray(t *testing.T) {
	h := history.Handler(history.New(5))
	req := httptest.NewRequest(http.MethodGet, "/history", nil)
	rw := httptest.NewRecorder()
	h(rw, req)

	var entries []history.Entry
	if err := json.NewDecoder(rw.Body).Decode(&entries); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if len(entries) != 0 {
		t.Errorf("expected empty array, got %d entries", len(entries))
	}
}
