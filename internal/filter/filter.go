package filter

import (
	"fmt"
	"time"

	"github.com/user/logslice/internal/parser"
)

// Options holds filtering criteria for log entries.
type Options struct {
	TimeFrom  *time.Time
	TimeTo    *time.Time
	FieldKey  string
	FieldVal  string
	TimeField string
}

// Apply filters a slice of log entries according to the given Options.
// Entries are included only if they match all specified criteria.
// Time filtering uses the field named by TimeField (default: "time").
// Entries missing the time field or with unparseable timestamps are excluded
// when time filtering is active.
func Apply(entries []parser.Entry, opts Options) ([]parser.Entry, error) {
	timeField := opts.TimeField
	if timeField == "" {
		timeField = "time"
	}

	var result []parser.Entry
	for _, e := range entries {
		if opts.TimeFrom != nil || opts.TimeTo != nil {
			raw, ok := e.Fields[timeField]
			if !ok {
				continue
			}
			ts, err := parseTime(fmt.Sprintf("%v", raw))
			if err != nil {
				continue
			}
			if opts.TimeFrom != nil && ts.Before(*opts.TimeFrom) {
				continue
			}
			if opts.TimeTo != nil && ts.After(*opts.TimeTo) {
				continue
			}
		}

		if opts.FieldKey != "" {
			val, ok := e.Fields[opts.FieldKey]
			if !ok || fmt.Sprintf("%v", val) != opts.FieldVal {
				continue
			}
		}

		result = append(result, e)
	}
	return result, nil
}

// parseTime attempts to parse s using several common timestamp formats.
func parseTime(s string) (time.Time, error) {
	formats := []string{
		time.RFC3339,
		time.RFC3339Nano,
		"2006-01-02T15:04:05",
		"2006-01-02 15:04:05",
	}
	for _, f := range formats {
		if t, err := time.Parse(f, s); err == nil {
			return t, nil
		}
	}
	return time.Time{}, fmt.Errorf("unrecognized time format: %q", s)
}
