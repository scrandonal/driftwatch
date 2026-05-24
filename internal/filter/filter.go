// Package filter provides path-matching utilities for deciding which
// file-system paths driftwatch should monitor.
package filter

import (
	"path/filepath"
	"strings"
)

// Filter decides whether a given file path should be watched.
type Filter struct {
	includeGlobs []string
	excludeGlobs []string
}

// New returns a Filter that accepts paths matching any of the includeGlobs
// and rejects paths matching any of the excludeGlobs.
// An empty includeGlobs slice means "accept everything" (before exclusions).
func New(includeGlobs, excludeGlobs []string) *Filter {
	return &Filter{
		includeGlobs: includeGlobs,
		excludeGlobs: excludeGlobs,
	}
}

// Allow returns true when path should be watched.
func (f *Filter) Allow(path string) bool {
	// Normalise separators so patterns work on every OS.
	path = filepath.ToSlash(strings.TrimSpace(path))

	if matchesAny(path, f.excludeGlobs) {
		return false
	}

	if len(f.includeGlobs) == 0 {
		return true
	}

	return matchesAny(path, f.includeGlobs)
}

// matchesAny reports whether path matches at least one of the supplied globs.
// Invalid patterns are silently skipped.
func matchesAny(path string, globs []string) bool {
	for _, g := range globs {
		matched, err := filepath.Match(g, path)
		if err != nil {
			continue
		}
		if matched {
			return true
		}
		// Also try matching just the base name so patterns like "*.tmp" work
		// against full paths.
		if matched, err = filepath.Match(g, filepath.Base(path)); err == nil && matched {
			return true
		}
	}
	return false
}
