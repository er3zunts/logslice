package pivot_test

import (
	"testing"

	"github.com/user/logslice/internal/parser"
	"github.com/user/logslice/internal/pivot"
)

func makeEntry(fields map[string]interface{}) parser.Entry {
	return parser.Entry{Fields: fields}
}

func TestByField_Basic(t *testing.T) {
	entries := []parser.Entry{
		makeEntry(map[string]interface{}{"level": "info"}),
		makeEntry(map[string]interface{}{"level": "error"}),
		makeEntry(map[string]interface{}{"level": "info"}),
		makeEntry(map[string]interface{}{"level": "info"}),
	}

	results := pivot.ByField(entries, "level", false)

	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
	if results[0].Value != "info" || results[0].Count != 3 {
		t.Errorf("expected info=3, got %s=%d", results[0].Value, results[0].Count)
	}
	if results[1].Value != "error" || results[1].Count != 1 {
		t.Errorf("expected error=1, got %s=%d", results[1].Value, results[1].Count)
	}
}

func TestByField_MissingFieldIncluded(t *testing.T) {
	entries := []parser.Entry{
		makeEntry(map[string]interface{}{"level": "info"}),
		makeEntry(map[string]interface{}{"msg": "no level"}),
	}

	results := pivot.ByField(entries, "level", false)

	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
}

func TestByField_SkipMissing(t *testing.T) {
	entries := []parser.Entry{
		makeEntry(map[string]interface{}{"level": "warn"}),
		makeEntry(map[string]interface{}{"msg": "no level"}),
	}

	results := pivot.ByField(entries, "level", true)

	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Value != "warn" {
		t.Errorf("unexpected value: %s", results[0].Value)
	}
}

func TestByField_Empty(t *testing.T) {
	results := pivot.ByField(nil, "level", false)
	if len(results) != 0 {
		t.Errorf("expected empty results, got %d", len(results))
	}
}

func TestToEntries(t *testing.T) {
	results := []pivot.Result{
		{Value: "info", Count: 5},
		{Value: "error", Count: 2},
	}

	entries := pivot.ToEntries(results, "level")

	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}
	if entries[0].Fields["level"] != "info" || entries[0].Fields["count"] != 5 {
		t.Errorf("unexpected first entry: %+v", entries[0].Fields)
	}
}
