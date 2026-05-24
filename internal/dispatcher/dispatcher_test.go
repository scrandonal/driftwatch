package dispatcher_test

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/example/driftwatch/internal/dispatcher"
	"github.com/example/driftwatch/internal/filter"
)

type mockSender struct {
	mu     sync.Mutex
	events []dispatcher.Event
	err    error
}

func (m *mockSender) Dispatch(_ context.Context, path, oldHash, newHash string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.events = append(m.events, dispatcher.Event{Path: path, OldHash: oldHash, NewHash: newHash})
	return m.err
}

func (m *mockSender) received() []dispatcher.Event {
	m.mu.Lock()
	defer m.mu.Unlock()
	out := make([]dispatcher.Event, len(m.events))
	copy(out, m.events)
	return out
}

func TestDispatcher_ForwardsEvent(t *testing.T) {
	sender := &mockSender{}
	d := dispatcher.New(dispatcher.Config{Sender: sender})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go d.Run(ctx)

	d.Send(dispatcher.Event{Path: "/etc/app.conf", OldHash: "aaa", NewHash: "bbb"})

	time.Sleep(50 * time.Millisecond)
	evts := sender.received()
	if len(evts) != 1 {
		t.Fatalf("expected 1 event, got %d", len(evts))
	}
	if evts[0].Path != "/etc/app.conf" {
		t.Errorf("unexpected path: %s", evts[0].Path)
	}
}

func TestDispatcher_FilterDropsEvent(t *testing.T) {
	sender := &mockSender{}
	f := filter.New(filter.Config{Excludes: []string{"*.tmp"}})
	d := dispatcher.New(dispatcher.Config{Sender: sender, Filter: f})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go d.Run(ctx)

	d.Send(dispatcher.Event{Path: "/tmp/work.tmp", OldHash: "aaa", NewHash: "bbb"})

	time.Sleep(50 * time.Millisecond)
	if len(sender.received()) != 0 {
		t.Fatal("expected event to be filtered out")
	}
}

func TestDispatcher_ContextCancellation(t *testing.T) {
	sender := &mockSender{}
	d := dispatcher.New(dispatcher.Config{Sender: sender})

	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan struct{})
	go func() {
		d.Run(ctx)
		close(done)
	}()

	cancel()
	select {
	case <-done:
	case <-time.After(500 * time.Millisecond):
		t.Fatal("Run did not return after context cancellation")
	}
}
