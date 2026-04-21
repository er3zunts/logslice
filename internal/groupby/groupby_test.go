package groupby

import (
	"testing"

	"github.com/logslice/logslice/internal/parser"
)

func makeEntry(fields map[string]interface{}) parser.Entry {
	return parser.Entry{Fields: fields}
}

func TestApply_BasicGrouping(t *testing.T) {
	entries := []parser.Entry{
		makeEntry(map[string]interface{}{"level": "info"}),
		makeEntry(map[string]interface{}{"level": "error"}),
		makeEntry(map[string]interface{}{"level": "info"}),
	}
	groups := Apply(entries, Options{Field: "level"})
	if len(groups) != 2 {
		t.Fatalf("expected 2 groups, got %d", len(groups))
	}
	// sorted: error, info
	if groups[0].Key != "error" || groups[0].Count != 1 {
		t.Errorf("unexpected first group: %+v", groups[0])
	}
	if groups[1].Key != "info" || groups[1].Count != 2 {
		t.Errorf("unexpected second group: %+v", groups[1])
	}
}

func TestApply_MissingField(t *testing.T) {
	entries := []parser.Entry{
		makeEntry(map[string]interface{}{"level": "info"}),
		makeEntry(map[string]interface{}{"msg": "hello"}),
	}
	groups := Apply(entries, Options{Field: "level"})
	if len(groups) != 2 {
		t.Fatalf("expected 2 groups, got %d", len(groups))
	}
	keys := map[string]bool{}
	for _, g := range groups {
		keys[g.Key] = true
	}
	if !keys["(missing)"] {
		t.Error("expected (missing) group")
	}
}

func TestApply_KeepRows(t *testing.T) {
	entries := []parser.Entry{
		makeEntry(map[string]interface{}{"level": "info", "msg": "a"}),
		makeEntry(map[string]interface{}{"level": "info", "msg": "b"}),
	}
	groups := Apply(entries, Options{Field: "level", KeepRows: true})
	if len(groups) != 1 {
		t.Fatalf("expected 1 group, got %d", len(groups))
	}
	if len(groups[0].Items) != 2 {
		t.Errorf("expected 2 items in group, got %d", len(groups[0].Items))
	}
}

func TestApply_Empty(t *testing.T) {
	groups := Apply([]parser.Entry{}, Options{Field: "level"})
	if len(groups) != 0 {
		t.Errorf("expected empty result, got %d groups", len(groups))
	}
}
