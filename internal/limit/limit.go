// Package limit provides functionality to cap the number of log entries
// returned, with support for offset-based pagination.
package limit

import "github.com/logslice/logslice/internal/parser"

// Options controls how the limit is applied.
type Options struct {
	// Max is the maximum number of entries to return. Zero means no limit.
	Max int
	// Offset is the number of entries to skip before applying Max.
	Offset int
}

// Apply returns a slice of entries from the input, skipping the first
// Offset entries and returning at most Max entries. If Max is zero,
// all entries after the offset are returned.
func Apply(entries []parser.Entry, opts Options) []parser.Entry {
	if opts.Offset < 0 {
		opts.Offset = 0
	}
	if opts.Max < 0 {
		opts.Max = 0
	}

	if opts.Offset >= len(entries) {
		return []parser.Entry{}
	}

	sliced := entries[opts.Offset:]

	if opts.Max == 0 || opts.Max >= len(sliced) {
		return sliced
	}

	return sliced[:opts.Max]
}

// Page is a convenience wrapper that computes offset from a 1-based page
// number and a page size, then calls Apply.
func Page(entries []parser.Entry, page, pageSize int) []parser.Entry {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 1
	}
	return Apply(entries, Options{
		Offset: (page - 1) * pageSize,
		Max:    pageSize,
	})
}
