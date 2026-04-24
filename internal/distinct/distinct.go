// Package distinct provides deduplication of log entries by unique field values.
package distinct

import "github.com/logslice/logslice/internal/parser"

// Options controls how distinct filtering is applied.
type Options struct {
	// Fields to consider when determining uniqueness.
	// If empty, all fields are used.
	Fields []string
	// Limit caps the number of distinct entries returned (0 = unlimited).
	Limit int
}

// key builds a comparable string key from the specified fields of an entry.
func key(entry parser.Entry, fields []string) string {
	if len(fields) == 0 {
		// Use all fields in sorted order via the entry's raw map.
		var b []byte
		for _, f := range sortedKeys(entry) {
			b = append(b, f...)
			b = append(b, '=')
			v := entry[f]
			b = append(b, fmt.Sprintf("%v", v)...)
			b = append(b, ';')
		}
		return string(b)
	}
	var b []byte
	for _, f := range fields {
		v, _ := entry[f]
		b = append(b, f...)
		b = append(b, '=')
		b = append(b, fmt.Sprintf("%v", v)...)
		b = append(b, ';')
	}
	return string(b)
}

// sortedKeys returns the keys of an entry in sorted order.
func sortedKeys(entry parser.Entry) []string {
	keys := make([]string, 0, len(entry))
	for k := range entry {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

// Apply returns only entries with distinct values for the given fields.
func Apply(entries []parser.Entry, opts Options) []parser.Entry {
	seen := make(map[string]struct{})
	result := make([]parser.Entry, 0, len(entries))
	for _, e := range entries {
		k := key(e, opts.Fields)
		if _, exists := seen[k]; exists {
			continue
		}
		seen[k] = struct{}{}
		result = append(result, e)
		if opts.Limit > 0 && len(result) >= opts.Limit {
			break
		}
	}
	return result
}
