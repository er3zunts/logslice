package window

import (
	"testing"
	"time"

	"github.com/nicholasgasior/logslice/internal/parser"
)

func makeEntry(ts string) parser.Entry {
	return parser.Entry{"time": ts, "msg": "hello"}
}

func TestApply_Empty(t *testing.T) {
	buckets, err := Apply(nil, "time", time.Minute, Tumbling)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(buckets) != 0 {
		t.Fatalf("expected 0 buckets, got %d", len(buckets))
	}
}

func TestApply_InvalidDuration(t *testing.T) {
	entries := []parser.Entry{makeEntry("2024-01-01T00:00:00Z")}
	_, err := Apply(entries, "time", 0, Tumbling)
	if err == nil {
		t.Fatal("expected error for zero duration")
	}
}

func TestApply_Tumbling_SingleBucket(t *testing.T) {
	entries := []parser.Entry{
		makeEntry("2024-01-01T00:00:10Z"),
		makeEntry("2024-01-01T00:00:30Z"),
		makeEntry("2024-01-01T00:00:50Z"),
	}
	buckets, err := Apply(entries, "time", time.Minute, Tumbling)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	total := 0
	for _, b := range buckets {
		total += len(b.Entries)
	}
	if total != 3 {
		t.Errorf("expected 3 total entries across buckets, got %d", total)
	}
}

func TestApply_Tumbling_MultipleBuckets(t *testing.T) {
	entries := []parser.Entry{
		makeEntry("2024-01-01T00:00:10Z"),
		makeEntry("2024-01-01T00:01:10Z"),
		makeEntry("2024-01-01T00:02:10Z"),
	}
	buckets, err := Apply(entries, "time", time.Minute, Tumbling)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(buckets) < 3 {
		t.Errorf("expected at least 3 buckets, got %d", len(buckets))
	}
}

func TestApply_Sliding_OverlapsEntries(t *testing.T) {
	entries := []parser.Entry{
		makeEntry("2024-01-01T00:00:10Z"),
		makeEntry("2024-01-01T00:00:40Z"),
		makeEntry("2024-01-01T00:01:10Z"),
	}
	buckets, err := Apply(entries, "time", time.Minute, Sliding)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// With sliding windows entries can appear in more than one bucket.
	total := 0
	for _, b := range buckets {
		total += len(b.Entries)
	}
	if total < 3 {
		t.Errorf("expected at least 3 entry appearances across sliding buckets, got %d", total)
	}
}

func TestApply_MissingTimestampField(t *testing.T) {
	entries := []parser.Entry{{"msg": "no timestamp"}}
	_, err := Apply(entries, "time", time.Minute, Tumbling)
	if err == nil {
		t.Fatal("expected error when no timestamps are parseable")
	}
}

func TestBucketLabel(t *testing.T) {
	start, _ := time.Parse(time.RFC3339, "2024-01-01T00:00:00Z")
	end := start.Add(time.Minute)
	b := Bucket{Start: start, End: end, Entries: []parser.Entry{makeEntry("2024-01-01T00:00:10Z")}}
	label := b.Label()
	if label == "" {
		t.Error("expected non-empty label")
	}
}
