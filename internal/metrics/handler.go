package metrics

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// responsePayload is the JSON shape returned by the metrics HTTP handler.
type responsePayload struct {
	UptimeSeconds float64 `json:"uptime_seconds"`
	FilesChecked  int64   `json:"files_checked"`
	ChangesFound  int64   `json:"changes_found"`
	AlertsSent    int64   `json:"alerts_sent"`
	AlertErrors   int64   `json:"alert_errors"`
	WatchErrors   int64   `json:"watch_errors"`
}

// Handler returns an http.HandlerFunc that serialises a Snapshot to JSON.
// It is intended to be mounted at a diagnostics endpoint such as /metrics.
func Handler(c *Collector) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		snap := c.Snapshot()
		payload := responsePayload{
			UptimeSeconds: snap.Uptime.Seconds(),
			FilesChecked:  snap.FilesChecked,
			ChangesFound:  snap.ChangesFound,
			AlertsSent:    snap.AlertsSent,
			AlertErrors:   snap.AlertErrors,
			WatchErrors:   snap.WatchErrors,
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(payload); err != nil {
			http.Error(w, fmt.Sprintf("encode error: %v", err), http.StatusInternalServerError)
		}
	}
}
