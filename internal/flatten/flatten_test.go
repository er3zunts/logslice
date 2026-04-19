package flatten

import (
	"testing"

	"github.com/logslice/logslice/internal/parser"
)

func makeEntry(fields map[string]interface{}) parser.Entry {
	return parser.Entry{Fields: fields}
}

func TestApply_NoNesting(t *testing.T) {
	entries := []parser.Entry{makeEntry(map[string]interface{}{"level": "info", "msg": "hello"})}
	out := Apply(entries, ".")
	if out[0].Fields["level"] != "info" {
		t.Errorf("expected info, got %v", out[0].Fields["level"])
	}
}

func TestApply_SingleNesting(t *testing.T) {
	fields := map[string]interface{}{
		"meta": map[string]interface{}{"user": "alice", "id": "42"},
	}
	out := Apply([]parser.Entry{makeEntry(fields)}, ".")
	if out[0].Fields["meta.user"] != "alice" {
		t.Errorf("expected alice, got %v", out[0].Fields["meta.user"])
	}
	if out[0].Fields["meta.id"] != "42" {
		t.Errorf("expected 42, got %v", out[0].Fields["meta.id"])
	}
	if _, ok := out[0].Fields["meta"]; ok {
		t.Error("expected nested key 'meta' to be removed")
	}
}

func TestApply_DeepNesting(t *testing.T) {
	fields := map[string]interface{}{
		"a": map[string]interface{}{
			"b": map[string]interface{}{
				"c": "deep",
			},
		},
	}
	out := Apply([]parser.Entry{makeEntry(fields)}, ".")
	if out[0].Fields["a.b.c"] != "deep" {
		t.Errorf("expected deep, got %v", out[0].Fields["a.b.c"])
	}
}

func TestApply_CustomSeparator(t *testing.T) {
	fields := map[string]interface{}{
		"x": map[string]interface{}{"y": "val"},
	}
	out := Apply([]parser.Entry{makeEntry(fields)}, "_")
	if out[0].Fields["x_y"] != "val" {
		t.Errorf("expected val, got %v", out[0].Fields["x_y"])
	}
}

func TestApply_EmptyEntries(t *testing.T) {
	out := Apply([]parser.Entry{}, ".")
	if len(out) != 0 {
		t.Errorf("expected empty, got %d", len(out))
	}
}
