package summarize

import (
	"testing"

	"github.com/user/logslice/internal/parser"
)

func makeEntry(fields map[string]interface{}) parser.Entry {
	return parser.Entry{Fields: fields}
}

func TestApply_Empty(t *testing.T) {
	results := Apply(nil, "service", "latency")
	if len(results) != 0 {
		t.Fatalf("expected 0 results, got %d", len(results))
	}
}

func TestApply_NoGroupField(t *testing.T) {
	entries := []parser.Entry{
		makeEntry(map[string]interface{}{"latency": float64(10)}),
		makeEntry(map[string]interface{}{"latency": float64(20)}),
		makeEntry(map[string]interface{}{"latency": float64(30)}),
	}
	results := Apply(entries, "", "latency")
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	r := results[0]
	if r.GroupKey != "(all)" {
		t.Errorf("expected group key '(all)', got %q", r.GroupKey)
	}
	if r.Stats.Count != 3 {
		t.Errorf("expected count 3, got %d", r.Stats.Count)
	}
	if r.Stats.Sum != 60 {
		t.Errorf("expected sum 60, got %v", r.Stats.Sum)
	}
	if r.Stats.Min != 10 {
		t.Errorf("expected min 10, got %v", r.Stats.Min)
	}
	if r.Stats.Max != 30 {
		t.Errorf("expected max 30, got %v", r.Stats.Max)
	}
	if r.Stats.Avg() != 20 {
		t.Errorf("expected avg 20, got %v", r.Stats.Avg())
	}
}

func TestApply_WithGroupField(t *testing.T) {
	entries := []parser.Entry{
		makeEntry(map[string]interface{}{"svc": "api", "dur": float64(5)}),
		makeEntry(map[string]interface{}{"svc": "api", "dur": float64(15)}),
		makeEntry(map[string]interface{}{"svc": "db", "dur": float64(50)}),
	}
	results := Apply(entries, "svc", "dur")
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
	// results are sorted by group key
	if results[0].GroupKey != "api" {
		t.Errorf("expected first group 'api', got %q", results[0].GroupKey)
	}
	if results[0].Stats.Count != 2 {
		t.Errorf("expected api count 2, got %d", results[0].Stats.Count)
	}
	if results[1].GroupKey != "db" {
		t.Errorf("expected second group 'db', got %q", results[1].GroupKey)
	}
	if results[1].Stats.Sum != 50 {
		t.Errorf("expected db sum 50, got %v", results[1].Stats.Sum)
	}
}

func TestApply_SkipsNonNumeric(t *testing.T) {
	entries := []parser.Entry{
		makeEntry(map[string]interface{}{"latency": "fast"}),
		makeEntry(map[string]interface{}{"latency": float64(99)}),
	}
	results := Apply(entries, "", "latency")
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Stats.Count != 1 {
		t.Errorf("expected count 1, got %d", results[0].Stats.Count)
	}
}

func TestApply_MissingGroupKey(t *testing.T) {
	entries := []parser.Entry{
		makeEntry(map[string]interface{}{"val": float64(7)}),
	}
	results := Apply(entries, "svc", "val")
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].GroupKey != "(missing)" {
		t.Errorf("expected group key '(missing)', got %q", results[0].GroupKey)
	}
}
