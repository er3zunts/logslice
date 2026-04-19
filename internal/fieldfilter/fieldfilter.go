// Package fieldfilter provides functionality to include or exclude specific
// fields from log entries.
package fieldfilter

import "github.com/logslice/logslice/internal/parser"

// Mode controls whether the field list is an allowlist or denylist.
type Mode int

const (
	Include Mode = iota
	Exclude
)

// Options configures field filtering behavior.
type Options struct {
	Fields []string
	Mode   Mode
}

// Apply filters fields from each entry according to the given options.
// In Include mode, only listed fields are kept.
// In Exclude mode, listed fields are removed.
func Apply(entries []parser.Entry, opts Options) []parser.Entry {
	if len(opts.Fields) == 0 {
		return entries
	}

	set := make(map[string]struct{}, len(opts.Fields))
	for _, f := range opts.Fields {
		set[f] = struct{}{}
	}

	result := make([]parser.Entry, 0, len(entries))
	for _, e := range entries {
		filtered := make(map[string]interface{})
		for k, v := range e.Fields {
			_, inSet := set[k]
			if opts.Mode == Include && inSet {
				filtered[k] = v
			} else if opts.Mode == Exclude && !inSet {
				filtered[k] = v
			}
		}
		result = append(result, parser.Entry{Fields: filtered})
	}
	return result
}
