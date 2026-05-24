// Package healthz provides a simple HTTP health-check endpoint for
// driftwatch. It reports whether the service is alive and exposes basic
// runtime information such as uptime and the number of files being watched.
package healthz

import (
	"encoding/json"
	"net/http"
	"time"
)

// StatusProvider is satisfied by any type that can report the current number
// of watched paths and the time the service started.
type StatusProvider interface {
	WatchedPaths() int
	StartedAt() time.Time
}

// Response is the JSON body returned by the health endpoint.
type Response struct {
	Status      string  `json:"status"`
	UptimeSeconds float64 `json:"uptime_seconds"`
	WatchedPaths int     `json:"watched_paths"`
}

// Handler returns an http.HandlerFunc that writes a JSON health response.
// It accepts only GET and HEAD requests; all other methods receive 405.
func Handler(sp StatusProvider) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet && r.Method != http.MethodHead {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		resp := Response{
			Status:        "ok",
			UptimeSeconds: time.Since(sp.StartedAt()).Seconds(),
			WatchedPaths:  sp.WatchedPaths(),
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(resp)
	}
}
