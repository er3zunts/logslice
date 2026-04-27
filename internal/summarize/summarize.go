// Package summarize reduces a stream of log entries to aggregate numeric summaries
// grouped by one or more fields (min, max, sum, avg, count).
package summarize

import (
	"fmt"
	"math"
	"sort"

	"github.com/user/logslice/internal/parser"
)

// Stats holds aggregate values for a numeric field.
type Stats struct {
	Count int64
	Sum   float64
	Min   float64
	Max   float64
}

// Avg returns the arithmetic mean, or 0 if Count is zero.
func (s Stats) Avg() float64 {
	if s.Count == 0 {
		return 0
	}
	return s.Sum / float64(s.Count)
}

// Result is a single summarized row.
type Result struct {
	GroupKey string
	Field    string
	Stats    Stats
}

// Apply groups entries by groupField and computes aggregate statistics for
// valueField. If groupField is empty every entry is placed in a single group.
func Apply(entries []parser.Entry, groupField, valueField string) []Result {
	type bucket struct {
		Stats
	}

	buckets := map[string]*bucket{}
	order := []string{}

	for _, e := range entries {
		groupKey := "(all)"
		if groupField != "" {
			if v, ok := e.Fields[groupField]; ok {
				groupKey = fmt.Sprintf("%v", v)
			} else {
				groupKey = "(missing)"
			}
		}

		v, ok := toFloat(e.Fields[valueField])
		if !ok {
			continue
		}

		b, exists := buckets[groupKey]
		if !exists {
			b = &bucket{Stats: Stats{Min: math.MaxFloat64, Max: -math.MaxFloat64}}
			buckets[groupKey] = b
			order = append(order, groupKey)
		}
		b.Count++
		b.Sum += v
		if v < b.Min {
			b.Min = v
		}
		if v > b.Max {
			b.Max = v
		}
	}

	sort.Strings(order)
	results := make([]Result, 0, len(order))
	for _, k := range order {
		results = append(results, Result{
			GroupKey: k,
			Field:    valueField,
			Stats:    buckets[k].Stats,
		})
	}
	return results
}

func toFloat(v interface{}) (float64, bool) {
	if v == nil {
		return 0, false
	}
	switch n := v.(type) {
	case float64:
		return n, true
	case int:
		return float64(n), true
	case int64:
		return float64(n), true
	}
	return 0, false
}
