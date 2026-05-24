// Package version exposes build-time version information for driftwatch.
package version

import "fmt"

// These variables are set at build time via -ldflags.
var (
	// Version is the semantic version string (e.g. "1.2.3").
	Version = "dev"

	// Commit is the short Git SHA of the build.
	Commit = "unknown"

	// Date is the UTC build timestamp in RFC3339 format.
	Date = "unknown"
)

// Info holds structured version metadata.
type Info struct {
	Version string `json:"version"`
	Commit  string `json:"commit"`
	Date    string `json:"date"`
}

// Get returns the current build Info.
func Get() Info {
	return Info{
		Version: Version,
		Commit:  Commit,
		Date:    Date,
	}
}

// String returns a human-readable version string.
func (i Info) String() string {
	return fmt.Sprintf("driftwatch %s (commit=%s, built=%s)", i.Version, i.Commit, i.Date)
}
