// Package split provides functionality to split a stream of log entries
// into multiple named buckets based on the value of a specified field.
package split

import (
	"fmt"

	"github.com/yourorg/logslice/internal/parser"
)

// Result holds the bucketed entries keyed by the split field's value.
type Result struct {
	Buckets map[string][]parser.Entry
	Order   []string // insertion-ordered list of bucket keys
}

// Apply partitions entries into buckets based on the value of field.
// Entries that are missing the field are placed in a bucket named by
// missingKey (use an empty string to discard them).
func Apply(entries []parser.Entry, field string, missingKey string) Result {
	r := Result{
		Buckets: make(map[string][]parser.Entry),
	}

	for _, e := range entries {
		val, ok := e[field]
		var key string
		if !ok || val == nil {
			if missingKey == "" {
				continue
			}
			key = missingKey
		} else {
			key = fmt.Sprintf("%v", val)
		}

		if _, exists := r.Buckets[key]; !exists {
			r.Order = append(r.Order, key)
		}
		r.Buckets[key] = append(r.Buckets[key], e)
	}

	return r
}

// Flatten returns all entries in bucket order, optionally injecting a
// new field (labelField) whose value is the bucket key. Pass an empty
// labelField to skip injection.
func Flatten(r Result, labelField string) []parser.Entry {
	var out []parser.Entry
	for _, key := range r.Order {
		for _, e := range r.Buckets[key] {
			if labelField != "" {
				copy := make(parser.Entry, len(e))
				for k, v := range e {
					copy[k] = v
				}
				copy[labelField] = key
				out = append(out, copy)
			} else {
				out = append(out, e)
			}
		}
	}
	return out
}
