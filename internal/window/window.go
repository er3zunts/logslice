// Package window provides sliding and tumbling time-window aggregation over log entries.
package window

import (
	"fmt"
	"time"

	"github.com/nicholasgasior/logslice/internal/parser"
)

// WindowType controls how windows are applied.
type WindowType string

const (
	Tumbling WindowType = "tumbling"
	Sliding  WindowType = "sliding"
)

// Bucket holds entries that fall within a single time window.
type Bucket struct {
	Start   time.Time
	End     time.Time
	Entries []parser.Entry
}

// Label returns a human-readable label for the bucket.
func (b Bucket) Label() string {
	return fmt.Sprintf("%s – %s (%d entries)", b.Start.Format(time.RFC3339), b.End.Format(time.RFC3339), len(b.Entries))
}

// Apply partitions entries into time buckets of the given duration.
// The timestamp field name and window type must be provided.
// For sliding windows the step is half the duration.
func Apply(entries []parser.Entry, field string, dur time.Duration, wtype WindowType) ([]Bucket, error) {
	if len(entries) == 0 {
		return nil, nil
	}
	if dur <= 0 {
		return nil, fmt.Errorf("window duration must be positive")
	}

	var buckets []Bucket
	step := dur
	if wtype == Sliding {
		step = dur / 2
		if step == 0 {
			step = dur
		}
	}

	// Determine overall time range.
	var minT, maxT time.Time
	for _, e := range entries {
		t, err := parseTimestamp(e, field)
		if err != nil {
			continue
		}
		if minT.IsZero() || t.Before(minT) {
			minT = t
		}
		if maxT.IsZero() || t.After(maxT) {
			maxT = t
		}
	}
	if minT.IsZero() {
		return nil, fmt.Errorf("no parseable timestamps found in field %q", field)
	}

	for start := minT.Truncate(step); !start.After(maxT); start = start.Add(step) {
		end := start.Add(dur)
		bucket := Bucket{Start: start, End: end}
		for _, e := range entries {
			t, err := parseTimestamp(e, field)
			if err != nil {
				continue
			}
			if !t.Before(start) && t.Before(end) {
				bucket.Entries = append(bucket.Entries, e)
			}
		}
		buckets = append(buckets, bucket)
	}
	return buckets, nil
}

func parseTimestamp(e parser.Entry, field string) (time.Time, error) {
	v, ok := e[field]
	if !ok {
		return time.Time{}, fmt.Errorf("field missing")
	}
	s, ok := v.(string)
	if !ok {
		return time.Time{}, fmt.Errorf("field not a string")
	}
	formats := []string{time.RFC3339Nano, time.RFC3339, "2006-01-02T15:04:05"}
	for _, f := range formats {
		if t, err := time.Parse(f, s); err == nil {
			return t, nil
		}
	}
	return time.Time{}, fmt.Errorf("unrecognised timestamp format: %s", s)
}
