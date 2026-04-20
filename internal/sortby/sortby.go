// Package sortby provides stable sorting of log entries by a named field.
package sortby

import (
	"sort"
	"strconv"

	"github.com/user/logslice/internal/parser"
)

// Order controls sort direction.
type Order int

const (
	Ascending  Order = iota
	Descending Order = iota
)

// Apply returns a new slice of entries sorted by field.
// Numeric fields are compared numerically; all others lexicographically.
// Entries missing the field are placed at the end.
func Apply(entries []parser.Entry, field string, order Order) []parser.Entry {
	out := make([]parser.Entry, len(entries))
	copy(out, entries)

	sort.SliceStable(out, func(i, j int) bool {
		vi, oki := out[i].Fields[field]
		vj, okj := out[j].Fields[field]

		// missing fields sink to the end
		if !oki && !okj {
			return false
		}
		if !oki {
			return false
		}
		if !okj {
			return true
		}

		less := compareValues(vi, vj)
		if order == Descending {
			return !less
		}
		return less
	})
	return out
}

// compareValues returns true when a < b, using numeric comparison when possible.
func compareValues(a, b interface{}) bool {
	af, aerr := toFloat(a)
	bf, berr := toFloat(b)
	if aerr == nil && berr == nil {
		return af < bf
	}
	return str(a) < str(b)
}

func toFloat(v interface{}) (float64, error) {
	switch t := v.(type) {
	case float64:
		return t, nil
	case int:
		return float64(t), nil
	case string:
		return strconv.ParseFloat(t, 64)
	}
	return 0, strconv.ErrSyntax
}

func str(v interface{}) string {
	if v == nil {
		return ""
	}
	switch t := v.(type) {
	case string:
		return t
	}
	return ""
}
