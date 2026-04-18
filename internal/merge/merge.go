// Package merge provides functionality to merge multiple sorted log streams.
package merge

import (
	"sort"
	"time"

	"github.com/user/logslice/internal/parser"
)

// TimeKey is the field used to extract timestamps for sorting.
const TimeKey = "time"

// ByTime merges multiple slices of log entries and returns them sorted by
// the timestamp field. Entries without a parseable timestamp are appended
// at the end in their original order.
func ByTime(streams ...[]parser.Entry) []parser.Entry {
	total := 0
	for _, s := range streams {
		total += len(s)
	}
	merged := make([]parser.Entry, 0, total)
	for _, s := range streams {
		merged = append(merged, s...)
	}

	timedFormats := []string{
		time.RFC3339Nano,
		time.RFC3339,
		"2006-01-02T15:04:05",
		"2006-01-02 15:04:05",
	}

	parseTS := func(e parser.Entry) (time.Time, bool) {
		v, ok := e.Fields[TimeKey]
		if !ok {
			return time.Time{}, false
		}
		s, ok := v.(string)
		if !ok {
			return time.Time{}, false
		}
		for _, f := range timedFormats {
			if t, err := time.Parse(f, s); err == nil {
				return t, true
			}
		}
		return time.Time{}, false
	}

	sort.SliceStable(merged, func(i, j int) bool {
		ti, oki := parseTS(merged[i])
		tj, okj := parseTS(merged[j])
		if oki && okj {
			return ti.Before(tj)
		}
		if oki {
			return true
		}
		return false
	})

	return merged
}
