// Package watcher monitors one or more files on disk for content changes.
//
// A Watcher polls each registered path at a configurable interval, computing
// a SHA-256 hash of the file contents on every tick. When the hash differs
// from the previously recorded value the Watcher emits a ChangeEvent on its
// Events channel so that callers can react (e.g. fire a webhook alert).
//
// # Basic usage
//
//	cfg := watcher.Config{
//		Paths:    []string{"/etc/app/config.yaml", "/etc/app/secrets.env"},
//		Interval: 10 * time.Second,
//		Logger:   slog.Default(),
//	}
//
//	w, err := watcher.New(cfg)
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	go func() {
//		for event := range w.Events {
//			fmt.Printf("changed: %s (old=%s new=%s)\n",
//				event.Path, event.OldHash, event.NewHash)
//		}
//	}()
//
//	ctx, cancel := context.WithCancel(context.Background())
//	defer cancel()
//	w.Run(ctx)
//
// The Watcher stops polling and closes the Events channel when the supplied
// context is cancelled. It is safe to range over Events after cancellation;
// the range will terminate once the channel is drained and closed.
//
// # Change events
//
// A ChangeEvent is emitted only when the hash actually changes between two
// consecutive polls. Transient read errors are logged and skipped; they do
// not produce events and do not stop the watcher.
package watcher
