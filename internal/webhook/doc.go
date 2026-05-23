// Package webhook provides a lightweight HTTP client for sending
// drift-detection notifications to a configured webhook endpoint.
//
// Usage:
//
//	client := webhook.New("https://example.com/hooks/driftwatch")
//	err := client.Send(webhook.Payload{
//		Event:    "file_changed",
//		FilePath: "/etc/app/config.yaml",
//		OldHash:  "abc123",
//		NewHash:  "def456",
//	})
//	if err != nil {
//		log.Printf("webhook notification failed: %v", err)
//	}
//
// The Payload is serialised as JSON and delivered via HTTP POST.
// A non-2xx response status is treated as an error.
package webhook
