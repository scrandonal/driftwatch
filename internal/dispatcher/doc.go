// Package dispatcher coordinates file-change events produced by the
// watcher and routes them to the notification pipeline.
//
// A Dispatcher applies an optional [filter.Filter] to decide whether
// an event should be forwarded, then calls the configured [Sender]
// (typically a [notify.Pipeline]) to deliver the alert.
//
// Events are buffered internally so that the watcher goroutine is
// never blocked waiting for network I/O.
//
// Typical usage:
//
//	d := dispatcher.New(dispatcher.Config{
//		Filter: f,
//		Sender: pipeline,
//		Log:    logger,
//	})
//	go d.Run(ctx)
//
//	// from watcher callback:
//	d.Send(dispatcher.Event{Path: p, OldHash: old, NewHash: new})
package dispatcher
