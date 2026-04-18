// Package dedupe provides deduplication of log entries based on field values.
package dedupe

import (
	"fmt"
	"strings"

	"github.com/user/logslice/internal/parser"
)

// Options controls deduplication behavior.
type Options struct {
	Fields []string // fields to use as dedup key; empty means all fields
	KeepFirst bool   // if true, keep first occurrence; otherwise keep last
}

// key builds a string key from the specified fields of an entry.
func key(entry parser.Entry, fields []string) string {
	if len(fields) == 0 {
		// use entire raw line as key
		return entry.Raw
	}
	parts := make([]string, 0, len(fields))
	for _, f := range fields {
		v, ok := entry.Fields[f]
		if !ok {
			parts = append(parts, "")
			continue
		}
		parts = append(parts, fmt.Sprintf("%v", v))
	}
	return strings.Join(parts, "\x00")
}

// Apply removes duplicate entries according to opts.
func Apply(entries []parser.Entry, opts Options) []parser.Entry {
	seen := make(map[string]int) // key -> index in result
	result := make([]parser.Entry, 0, len(entries))

	for _, e := range entries {
		k := key(e, opts.Fields)
		if idx, exists := seen[k]; exists {
			if !opts.KeepFirst {
				// replace with latest
				result[idx] = e
			}
			continue
		}
		seen[k] = len(result)
		result = append(result, e)
	}
	return result
}
