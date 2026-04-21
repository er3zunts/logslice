// Package groupby groups log entries by a field value and aggregates counts.
package groupby

import (
	"fmt"
	"sort"

	"github.com/logslice/logslice/internal/parser"
)

// Group holds aggregated data for a single field value.
type Group struct {
	Key   string
	Count int
	Items []parser.Entry
}

// Options controls how grouping is performed.
type Options struct {
	Field    string
	KeepRows bool // if true, retain individual entries in each group
}

// Apply groups entries by the given field and returns a slice of Groups
// sorted by key alphabetically.
func Apply(entries []parser.Entry, opts Options) []Group {
	index := make(map[string]*Group)
	order := []string{}

	for _, e := range entries {
		var key string
		if v, ok := e.Fields[opts.Field]; ok {
			key = fmt.Sprintf("%v", v)
		} else {
			key = "(missing)"
		}

		g, exists := index[key]
		if !exists {
			g = &Group{Key: key}
			index[key] = g
			order = append(order, key)
		}
		g.Count++
		if opts.KeepRows {
			g.Items = append(g.Items, e)
		}
	}

	sort.Strings(order)

	result := make([]Group, 0, len(order))
	for _, k := range order {
		result = append(result, *index[k])
	}
	return result
}
