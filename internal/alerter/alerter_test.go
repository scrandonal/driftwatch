package alerter_test

import (
	"context"
	"errors"
	"log"
	"os"
	"testing"
	"time"

	"github.com/user/driftwatch/internal/alerter"
)

// mockSender records calls to Send and can simulate errors.
type mockSender struct {
	called  bool
	payload any
	errToReturn error
}

func (m *mockSender) Send(_ context.Context, payload any) error {
	m.called = true
	m.payload = payload
	return m.errToReturn
}

func newTestAlerter(sender *mockSender) *alerter.Alerter {
	silent := log.New(os.Discard, "", 0)
	return alerter.New(sender, silent)
}

func TestAlerter_Notify_Success(t *testing.T) {
	sender := &mockSender{}
	a := newTestAlerter(sender)

	event := alerter.Event{
		FilePath:   "/etc/app/config.yaml",
		OldHash:    "abc123",
		NewHash:    "def456",
		DetectedAt: time.Now(),
	}

	if err := a.Notify(context.Background(), event); err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if !sender.called {
		t.Fatal("expected sender.Send to be called")
	}
	got, ok := sender.payload.(alerter.Event)
	if !ok {
		t.Fatal("expected payload to be of type alerter.Event")
	}
	if got.FilePath != event.FilePath {
		t.Errorf("expected FilePath %q, got %q", event.FilePath, got.FilePath)
	}
}

func TestAlerter_Notify_SenderError(t *testing.T) {
	sendErr := errors.New("connection refused")
	sender := &mockSender{errToReturn: sendErr}
	a := newTestAlerter(sender)

	event := alerter.Event{
		FilePath:   "/etc/app/config.yaml",
		OldHash:    "aaa",
		NewHash:    "bbb",
		DetectedAt: time.Now(),
	}

	err := a.Notify(context.Background(), event)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, sendErr) {
		t.Errorf("expected wrapped sendErr, got: %v", err)
	}
}

func TestAlerter_New_NilLogger(t *testing.T) {
	sender := &mockSender{}
	// Should not panic when logger is nil.
	a := alerter.New(sender, nil)
	if a == nil {
		t.Fatal("expected non-nil Alerter")
	}
}
