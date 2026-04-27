// Package reorder provides functionality to reorder fields in log entries.
package reorder

import "github.com/logslice/logslice/internal/parser"

// Apply reorders the fields in each entry so that the specified fields appear
// first (in the given order), followed by any remaining fields in their
// original iteration order. Fields listed that do not exist in an entry are
// silently skipped.
func Apply(entries []parser.Entry, fields []string) []parser.Entry {
	if len(entries) == 0 || len(fields) == 0 {
		return entries
	}

	result := make([]parser.Entry, 0, len(entries))
	for _, entry := range entries {
		result = append(result, reorderEntry(entry, fields))
	}
	return result
}

func reorderEntry(entry parser.Entry, fields []string) parser.Entry {
	out := make(parser.Entry, len(entry))

	// Track which keys have been placed.
	placed := make(map[string]bool, len(fields))

	// First, insert the requested fields in order.
	for _, f := range fields {
		if v, ok := entry[f]; ok {
			out[f] = v
			placed[f] = true
		}
	}

	// Then append remaining fields.
	for k, v := range entry {
		if !placed[k] {
			out[k] = v
		}
	}

	return out
}
