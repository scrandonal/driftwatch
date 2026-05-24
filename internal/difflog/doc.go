// Package difflog provides a lightweight, append-only log of file-change
// events detected by driftwatch.
//
// Each change is recorded as a newline-delimited JSON object containing:
//
//   - timestamp  – UTC time the change was observed
//   - path       – absolute path of the changed file
//   - prev_hash  – SHA-256 hex digest of the file before the change
//   - new_hash   – SHA-256 hex digest of the file after the change
//
// Usage:
//
//	f, err := os.OpenFile("drift.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
//	if err != nil { /* handle */ }
//	defer f.Close()
//
//	log, err := difflog.New(f)
//	if err != nil { /* handle */ }
//
//	log.Record(difflog.Entry{
//		Path:     "/etc/app/config.yaml",
//		PrevHash: prevSHA,
//		NewHash:  newSHA,
//	})
package difflog
