package webhook_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/user/driftwatch/internal/webhook"
)

func TestClient_Send_Success(t *testing.T) {
	var received webhook.Payload

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if ct := r.Header.Get("Content-Type"); ct != "application/json" {
			t.Errorf("expected application/json, got %s", ct)
		}
		if err := json.NewDecoder(r.Body).Decode(&received); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := webhook.New(server.URL)
	payload := webhook.Payload{
		Event:     "file_changed",
		FilePath:  "/etc/app/config.yaml",
		Timestamp: time.Now().UTC(),
		OldHash:   "abc123",
		NewHash:   "def456",
	}

	if err := client.Send(payload); err != nil {
		t.Fatalf("Send returned error: %v", err)
	}

	if received.Event != payload.Event {
		t.Errorf("event: got %q, want %q", received.Event, payload.Event)
	}
	if received.FilePath != payload.FilePath {
		t.Errorf("file_path: got %q, want %q", received.FilePath, payload.FilePath)
	}
	if received.OldHash != payload.OldHash {
		t.Errorf("old_hash: got %q, want %q", received.OldHash, payload.OldHash)
	}
}

func TestClient_Send_NonSuccessStatus(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	client := webhook.New(server.URL)
	err := client.Send(webhook.Payload{Event: "file_changed", FilePath: "/tmp/test.conf"})
	if err == nil {
		t.Fatal("expected error for non-2xx status, got nil")
	}
}

func TestClient_Send_InvalidURL(t *testing.T) {
	client := webhook.New("http://127.0.0.1:0/nonexistent")
	err := client.Send(webhook.Payload{Event: "file_changed", FilePath: "/tmp/test.conf"})
	if err == nil {
		t.Fatal("expected error for unreachable URL, got nil")
	}
}
