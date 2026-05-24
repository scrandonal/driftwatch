// Package healthz exposes a lightweight HTTP liveness endpoint for the
// driftwatch service.
//
// Usage:
//
//	http.Handle("/healthz", healthz.Handler(myStatusProvider))
//
// The endpoint accepts GET and HEAD requests and responds with a JSON
// document containing:
//
//	{
//	  "status": "ok",
//	  "uptime_seconds": 42.7,
//	  "watched_paths": 3
//	}
//
// Any other HTTP method results in a 405 Method Not Allowed response.
package healthz
