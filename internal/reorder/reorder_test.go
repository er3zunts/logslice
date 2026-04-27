package reorder

import (
	"testing"

	"github.com/logslice/logslice/internal/parser"
)

func makeEntry(kv ...interface{}) parser.Entry {
	e := make(parser.Entry)
	for i := 0; i+1 < len(kv); i += 2 {
		e[kv[i].(string)] = kv[i+1]
	}
	return e
}

func TestApply_NoFields(t *testing.T) {
	entries := []parser.Entry{makeEntry("a", 1, "b", 2)}
	out := Apply(entries, nil)
	if len(out) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(out))
	}
}

func TestApply_Empty(t *testing.T) {
	out := Apply(nil, []string{"a"})
	if len(out) != 0 {
		t.Fatalf("expected 0 entries, got %d", len(out))
	}
}

func TestApply_BasicReorder(t *testing.T) {
	entries := []parser.Entry{
		makeEntry("level", "info", "msg", "hello", "ts", "2024-01-01"),
	}
	out := Apply(entries, []string{"ts", "level", "msg"})
	if len(out) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(out))
	}
	e := out[0]
	if e["ts"] != "2024-01-01" || e["level"] != "info" || e["msg"] != "hello" {
		t.Errorf("unexpected entry contents: %v", e)
	}
}

func TestApply_MissingFieldSkipped(t *testing.T) {
	entries := []parser.Entry{
		makeEntry("level", "warn", "msg", "oops"),
	}
	out := Apply(entries, []string{"ts", "level", "msg"})
	if _, ok := out[0]["ts"]; ok {
		t.Errorf("expected missing field 'ts' to be absent")
	}
	if out[0]["level"] != "warn" {
		t.Errorf("expected level=warn, got %v", out[0]["level"])
	}
}

func TestApply_PreservesExtraFields(t *testing.T) {
	entries := []parser.Entry{
		makeEntry("level", "debug", "msg", "hi", "caller", "main.go"),
	}
	out := Apply(entries, []string{"level"})
	if out[0]["msg"] != "hi" {
		t.Errorf("expected extra field msg=hi to be preserved")
	}
	if out[0]["caller"] != "main.go" {
		t.Errorf("expected extra field caller=main.go to be preserved")
	}
}

func TestApply_MultipleEntries(t *testing.T) {
	entries := []parser.Entry{
		makeEntry("a", 1, "b", 2),
		makeEntry("a", 3, "b", 4),
	}
	out := Apply(entries, []string{"b", "a"})
	if len(out) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(out))
	}
	for _, e := range out {
		if _, ok := e["a"]; !ok {
			t.Errorf("expected field 'a' in entry")
		}
		if _, ok := e["b"]; !ok {
			t.Errorf("expected field 'b' in entry")
		}
	}
}
