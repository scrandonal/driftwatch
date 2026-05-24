package dispatcher

import (
	"time"

	"github.com/example/driftwatch/internal/debounce"
	"github.com/example/driftwatch/internal/filter"
)

// BuildConfig holds the raw values parsed from the application config
// and is used to construct a Dispatcher with sensible defaults.
type BuildConfig struct {
	// DebounceDuration is how long to wait after the last event before
	// forwarding. Zero disables debouncing.
	DebounceDuration time.Duration

	// IncludeGlobs and ExcludeGlobs are forwarded to the filter.
	IncludeGlobs []string
	ExcludeGlobs []string
}

// DefaultBuildConfig returns a BuildConfig with sensible defaults.
func DefaultBuildConfig() BuildConfig {
	return BuildConfig{
		DebounceDuration: 200 * time.Millisecond,
	}
}

// Build constructs a Config suitable for New from a BuildConfig and a
// Sender. The caller is responsible for providing a non-nil Sender.
func Build(bc BuildConfig, sender Sender) Config {
	var f *filter.Filter
	if len(bc.IncludeGlobs) > 0 || len(bc.ExcludeGlobs) > 0 {
		f = filter.New(filter.Config{
			Includes: bc.IncludeGlobs,
			Excludes: bc.ExcludeGlobs,
		})
	}

	var db *debounce.Debouncer
	if bc.DebounceDuration > 0 {
		db = debounce.New(bc.DebounceDuration)
	}

	return Config{
		Filter:   f,
		Debounce: db,
		Sender:   sender,
	}
}
