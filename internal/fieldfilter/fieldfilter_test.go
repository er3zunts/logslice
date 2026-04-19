package fieldfilter_test

import (
	"testing"

	"github.com/logslice/logslice/internal/fieldfilter"
	"github.com/logslice/logslice/internal/parser"
)

func makeEntry(fields map[string]interface{}) parser.Entry {
	return parser.Entry{Fields: fields}
}

func TestApply_NoFields(t *testing.T) {
	entries := []parser.Entry{makeEntry(map[string]interface{}{"a": 1, "b": 2})}
	out := fieldfilter.Apply(entries, fieldfilter.Options{})
	if len(out) != 1 || len(out[0].Fields) != 2 {
		t.Fatal("expected unchanged entries when no fields specified")
	}
}

func TestApply_IncludeMode(t *testing.T) {
	entries := []parser.Entry{
		makeEntry(map[string]interface{}{"level": "info", "msg": "hello", "ts": "now"}),
	}
	out := fieldfilter.Apply(entries, fieldfilter.Options{
		Fields: []string{"level", "msg"},
		Mode:   fieldfilter.Include,
	})
	if len(out) != 1 {
		t.Fatal("expected 1 entry")
	}
	if _, ok := out[0].Fields["ts"]; ok {
		t.Error("expected 'ts' to be excluded")
	}
	if _, ok := out[0].Fields["level"]; !ok {
		t.Error("expected 'level' to be included")
	}
}

func TestApply_ExcludeMode(t *testing.T) {
	entries := []parser.Entry{
		makeEntry(map[string]interface{}{"level": "info", "msg": "hello", "ts": "now"}),
	}
	out := fieldfilter.Apply(entries, fieldfilter.Options{
		Fields: []string{"ts"},
		Mode:   fieldfilter.Exclude,
	})
	if _, ok := out[0].Fields["ts"]; ok {
		t.Error("expected 'ts' to be removed")
	}
	if _, ok := out[0].Fields["msg"]; !ok {
		t.Error("expected 'msg' to remain")
	}
}

func TestApply_MultipleEntries(t *testing.T) {
	entries := []parser.Entry{
		makeEntry(map[string]interface{}{"a": 1, "b": 2}),
		makeEntry(map[string]interface{}{"a": 3, "b": 4}),
	}
	out := fieldfilter.Apply(entries, fieldfilter.Options{
		Fields: []string{"a"},
		Mode:   fieldfilter.Include,
	})
	if len(out) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(out))
	}
	for _, e := range out {
		if _, ok := e.Fields["b"]; ok {
			t.Error("expected 'b' to be excluded")
		}
	}
}
