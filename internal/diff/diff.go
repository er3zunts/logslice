// Package diff provides utilities for comparing consecutive log entries
// and emitting only those where a specified field value has changed.
package diff

import (
	"fmt"

	"github.com/user/logslice/internal/parser"
)

// Options controls how the diff operation behaves.
type Options struct {
	// Fields is the list of fields to watch for changes.
	// If empty, all fields are compared.
	Fields []string

	// EmitFirst controls whether the very first entry is always emitted
	// regardless of comparison (there is nothing to compare it against).
	EmitFirst bool

	// AnnotateField, if non-empty, adds a field to each emitted entry
	// whose value is the name of the field that changed.
	AnnotateField string
}

// Apply returns only entries where at least one watched field differs
// from the same field in the previous entry. Order is preserved.
func Apply(entries []parser.Entry, opts Options) []parser.Entry {
	if len(entries) == 0 {
		return nil
	}

	var result []parser.Entry
	var prev parser.Entry

	for i, entry := range entries {
		if i == 0 {
			if opts.EmitFirst {
				result = append(result, entry)
			}
			prev = entry
			continue
		}

		changed, changedField := hasChanged(prev, entry, opts.Fields)
		if changed {
			out := entry
			if opts.AnnotateField != "" && changedField != "" {
				out = copyEntry(entry)
				out.Fields[opts.AnnotateField] = changedField
			}
			result = append(result, out)
		}
		prev = entry
	}

	return result
}

// hasChanged reports whether any watched field differs between a and b.
// It returns the name of the first changed field (or empty string if
// multiple fields are watched and several changed).
func hasChanged(a, b parser.Entry, fields []string) (bool, string) {
	watch := fields
	if len(watch) == 0 {
		// build union of all keys
		seen := make(map[string]struct{})
		for k := range a.Fields {
			seen[k] = struct{}{}
		}
		for k := range b.Fields {
			seen[k] = struct{}{}
		}
		for k := range seen {
			watch = append(watch, k)
		}
	}

	for _, f := range watch {
		av := fmt.Sprintf("%v", a.Fields[f])
		bv := fmt.Sprintf("%v", b.Fields[f])
		if av != bv {
			return true, f
		}
	}
	return false, ""
}

// copyEntry returns a shallow copy of e with a new Fields map.
func copyEntry(e parser.Entry) parser.Entry {
	out := parser.Entry{
		Raw:    e.Raw,
		Fields: make(map[string]interface{}, len(e.Fields)+1),
	}
	for k, v := range e.Fields {
		out.Fields[k] = v
	}
	return out
}
