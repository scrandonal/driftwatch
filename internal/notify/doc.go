// Package notify implements a notification pipeline that combines
// rate-limiting and configurable retry logic before handing off
// a payload to an underlying Sender (typically the webhook client).
//
// Usage:
//
//	cfg := notify.DefaultConfig()
//	cfg.MaxPerWindow = 5
//
//	p, err := notify.New(webhookClient, cfg, logger)
//	if err != nil { ... }
//
//	if err := p.Dispatch(ctx, payload); err != nil {
//	    // rate-limited or all retries exhausted
//	}
//
// The Pipeline is safe for concurrent use.
package notify
