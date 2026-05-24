package history

import (
	"encoding/json"
	"net/http"
	"strconv"
)

const defaultLimit = 20

// Handler returns an http.HandlerFunc that serves the recent change history
// as a JSON array. An optional `?limit=N` query parameter controls how many
// entries are returned (max 100).
func Handler(r *Ring) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		limit := defaultLimit
		if raw := req.URL.Query().Get("limit"); raw != "" {
			if n, err := strconv.Atoi(raw); err == nil && n > 0 {
				limit = n
			}
		}
		if limit > 100 {
			limit = 100
		}

		entries := r.Latest(limit)
		if entries == nil {
			entries = []Entry{}
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(entries)
	}
}
