// Package join provides functionality to enrich log entries by joining
// them with a static lookup table loaded from a JSON file.
package join

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/user/logslice/internal/parser"
)

// Options controls how the join is performed.
type Options struct {
	// LookupFile is the path to a JSON file containing an array of objects
	// used as the lookup table.
	LookupFile string

	// OnField is the field name present in both the log entry and the lookup
	// table that is used as the join key.
	OnField string

	// Prefix is prepended to every field name copied from the lookup row into
	// the log entry. An empty prefix copies fields as-is, potentially
	// overwriting existing fields.
	Prefix string

	// DropUnmatched discards log entries that have no matching row in the
	// lookup table when set to true. By default unmatched entries are passed
	// through unchanged.
	DropUnmatched bool
}

// loadLookup reads the lookup file and indexes rows by the value of onField.
func loadLookup(path, onField string) (map[string]parser.Entry, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("join: open lookup file: %w", err)
	}
	defer f.Close()

	var rows []parser.Entry
	if err := json.NewDecoder(f).Decode(&rows); err != nil {
		return nil, fmt.Errorf("join: decode lookup file: %w", err)
	}

	index := make(map[string]parser.Entry, len(rows))
	for _, row := range rows {
		v, ok := row[onField]
		if !ok {
			continue
		}
		key := fmt.Sprintf("%v", v)
		index[key] = row
	}
	return index, nil
}

// Apply performs a left join of entries against the lookup table described by
// opts. Each entry whose onField value matches a lookup row is enriched with
// the lookup row's fields (optionally prefixed). Entries without a match are
// either passed through or dropped depending on opts.DropUnmatched.
func Apply(entries []parser.Entry, opts Options) ([]parser.Entry, error) {
	if opts.OnField == "" {
		return nil, fmt.Errorf("join: OnField must not be empty")
	}

	lookup, err := loadLookup(opts.LookupFile, opts.OnField)
	if err != nil {
		return nil, err
	}

	out := make([]parser.Entry, 0, len(entries))
	for _, entry := range entries {
		v, ok := entry[opts.OnField]
		if !ok {
			if !opts.DropUnmatched {
				out = append(out, entry)
			}
			continue
		}

		key := fmt.Sprintf("%v", v)
		row, matched := lookup[key]
		if !matched {
			if !opts.DropUnmatched {
				out = append(out, entry)
			}
			continue
		}

		// Copy entry so we don't mutate the original slice.
		enriched := make(parser.Entry, len(entry))
		for k, val := range entry {
			enriched[k] = val
		}
		for k, val := range row {
			if k == opts.OnField {
				continue // don't duplicate the join key
			}
			enriched[opts.Prefix+k] = val
		}
		out = append(out, enriched)
	}
	return out, nil
}
