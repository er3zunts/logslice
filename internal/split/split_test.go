package split

import (
	"testing"

	"github.com/yourorg/logslice/internal/parser"
)

func makeEntry(fields map[string]interface{}) parser.Entry {
	e := make(parser.Entry)
	for k, v := range fields {
		e[k] = v
	}
	return e
}

func TestApply_BasicSplit(t *testing.T) {
	entries := []parser.Entry{
		makeEntry(map[string]interface{}{"level": "info", "msg": "a"}),
		makeEntry(map[string]interface{}{"level": "error", "msg": "b"}),
		makeEntry(map[string]interface{}{"level": "info", "msg": "c"}),
	}
	r := Apply(entries, "level", "")
	if len(r.Buckets["info"]) != 2 {
		t.Errorf("expected 2 info entries, got %d", len(r.Buckets["info"]))
	}
	if len(r.Buckets["error"]) != 1 {
		t.Errorf("expected 1 error entry, got %d", len(r.Buckets["error"]))
	}
}

func TestApply_MissingFieldDiscarded(t *testing.T) {
	entries := []parser.Entry{
		makeEntry(map[string]interface{}{"level": "info"}),
		makeEntry(map[string]interface{}{"msg": "no level"}),
	}
	r := Apply(entries, "level", "")
	if _, ok := r.Buckets[""]; ok {
		t.Error("expected missing-field entries to be discarded")
	}
	if len(r.Buckets["info"]) != 1 {
		t.Errorf("expected 1 info entry, got %d", len(r.Buckets["info"]))
	}
}

func TestApply_MissingFieldBucket(t *testing.T) {
	entries := []parser.Entry{
		makeEntry(map[string]interface{}{"level": "warn"}),
		makeEntry(map[string]interface{}{"msg": "no level"}),
	}
	r := Apply(entries, "level", "unknown")
	if len(r.Buckets["unknown"]) != 1 {
		t.Errorf("expected 1 unknown entry, got %d", len(r.Buckets["unknown"]))
	}
}

func TestApply_Empty(t *testing.T) {
	r := Apply([]parser.Entry{}, "level", "")
	if len(r.Buckets) != 0 {
		t.Error("expected empty result for empty input")
	}
}

func TestApply_OrderPreserved(t *testing.T) {
	entries := []parser.Entry{
		makeEntry(map[string]interface{}{"level": "debug"}),
		makeEntry(map[string]interface{}{"level": "info"}),
		makeEntry(map[string]interface{}{"level": "error"}),
	}
	r := Apply(entries, "level", "")
	expected := []string{"debug", "info", "error"}
	for i, k := range r.Order {
		if k != expected[i] {
			t.Errorf("order[%d]: expected %q, got %q", i, expected[i], k)
		}
	}
}

func TestFlatten_WithLabel(t *testing.T) {
	entries := []parser.Entry{
		makeEntry(map[string]interface{}{"level": "info", "msg": "hello"}),
	}
	r := Apply(entries, "level", "")
	out := Flatten(r, "_bucket")
	if len(out) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(out))
	}
	if out[0]["_bucket"] != "info" {
		t.Errorf("expected _bucket=info, got %v", out[0]["_bucket"])
	}
	if out[0]["msg"] != "hello" {
		t.Errorf("original field should be preserved")
	}
}

func TestFlatten_NoLabel(t *testing.T) {
	entries := []parser.Entry{
		makeEntry(map[string]interface{}{"level": "info"}),
	}
	r := Apply(entries, "level", "")
	out := Flatten(r, "")
	if _, ok := out[0]["_bucket"]; ok {
		t.Error("expected no label field when labelField is empty")
	}
}
