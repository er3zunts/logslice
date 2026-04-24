package distinct

import (
	"testing"

	"github.com/logslice/logslice/internal/parser"
)

func makeEntry(fields map[string]interface{}) parser.Entry {
	e := make(parser.Entry)
	for k, v := range fields {
		e[k] = v
	}
	return e
}

func TestApply_NoDuplicates(t *testing.T) {
	entries := []parser.Entry{
		makeEntry(map[string]interface{}{"level": "info", "msg": "a"}),
		makeEntry(map[string]interface{}{"level": "warn", "msg": "b"}),
		makeEntry(map[string]interface{}{"level": "error", "msg": "c"}),
	}
	out := Apply(entries, Options{Fields: []string{"level"}})
	if len(out) != 3 {
		t.Fatalf("expected 3, got %d", len(out))
	}
}

func TestApply_RemovesDuplicates(t *testing.T) {
	entries := []parser.Entry{
		makeEntry(map[string]interface{}{"level": "info", "msg": "first"}),
		makeEntry(map[string]interface{}{"level": "info", "msg": "second"}),
		makeEntry(map[string]interface{}{"level": "warn", "msg": "third"}),
	}
	out := Apply(entries, Options{Fields: []string{"level"}})
	if len(out) != 2 {
		t.Fatalf("expected 2, got %d", len(out))
	}
	if out[0]["msg"] != "first" {
		t.Errorf("expected first entry to be kept, got %v", out[0]["msg"])
	}
}

func TestApply_MultipleFields(t *testing.T) {
	entries := []parser.Entry{
		makeEntry(map[string]interface{}{"svc": "api", "level": "info"}),
		makeEntry(map[string]interface{}{"svc": "api", "level": "info"}),
		makeEntry(map[string]interface{}{"svc": "api", "level": "warn"}),
		makeEntry(map[string]interface{}{"svc": "db", "level": "info"}),
	}
	out := Apply(entries, Options{Fields: []string{"svc", "level"}})
	if len(out) != 3 {
		t.Fatalf("expected 3, got %d", len(out))
	}
}

func TestApply_WithLimit(t *testing.T) {
	entries := []parser.Entry{
		makeEntry(map[string]interface{}{"level": "info"}),
		makeEntry(map[string]interface{}{"level": "warn"}),
		makeEntry(map[string]interface{}{"level": "error"}),
	}
	out := Apply(entries, Options{Fields: []string{"level"}, Limit: 2})
	if len(out) != 2 {
		t.Fatalf("expected 2, got %d", len(out))
	}
}

func TestApply_Empty(t *testing.T) {
	out := Apply([]parser.Entry{}, Options{Fields: []string{"level"}})
	if len(out) != 0 {
		t.Fatalf("expected 0, got %d", len(out))
	}
}
