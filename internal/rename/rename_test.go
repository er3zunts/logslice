package rename_test

import (
	"testing"

	"github.com/logslice/logslice/internal/parser"
	"github.com/logslice/logslice/internal/rename"
)

func makeEntry(fields map[string]interface{}) parser.Entry {
	return parser.Entry{Fields: fields}
}

func TestApply_NoRules(t *testing.T) {
	entries := []parser.Entry{makeEntry(map[string]interface{}{"level": "info"})}
	out := rename.Apply(entries, nil)
	if len(out) != 1 || out[0].Fields["level"] != "info" {
		t.Fatal("expected unchanged entry")
	}
}

func TestApply_BasicRename(t *testing.T) {
	entries := []parser.Entry{makeEntry(map[string]interface{}{"msg": "hello", "ts": "2024-01-01"})}
	rules := []rename.Rule{{From: "msg", To: "message"}, {From: "ts", To: "timestamp"}}
	out := rename.Apply(entries, rules)
	if _, ok := out[0].Fields["msg"]; ok {
		t.Error("old field 'msg' should not exist")
	}
	if out[0].Fields["message"] != "hello" {
		t.Errorf("expected 'hello', got %v", out[0].Fields["message"])
	}
	if out[0].Fields["timestamp"] != "2024-01-01" {
		t.Errorf("expected '2024-01-01', got %v", out[0].Fields["timestamp"])
	}
}

func TestApply_MissingSourceField(t *testing.T) {
	entries := []parser.Entry{makeEntry(map[string]interface{}{"level": "warn"})}
	rules := []rename.Rule{{From: "nonexistent", To: "other"}}
	out := rename.Apply(entries, rules)
	if _, ok := out[0].Fields["other"]; ok {
		t.Error("destination field should not be created when source is missing")
	}
	if out[0].Fields["level"] != "warn" {
		t.Error("existing fields should be preserved")
	}
}

func TestApply_OverwritesDestination(t *testing.T) {
	entries := []parser.Entry{makeEntry(map[string]interface{}{"old": "value", "new": "existing"})}
	rules := []rename.Rule{{From: "old", To: "new"}}
	out := rename.Apply(entries, rules)
	if out[0].Fields["new"] != "value" {
		t.Errorf("expected destination to be overwritten with 'value', got %v", out[0].Fields["new"])
	}
	if _, ok := out[0].Fields["old"]; ok {
		t.Error("source field should be removed")
	}
}

func TestApply_MultipleEntries(t *testing.T) {
	entries := []parser.Entry{
		makeEntry(map[string]interface{}{"a": 1}),
		makeEntry(map[string]interface{}{"a": 2}),
	}
	rules := []rename.Rule{{From: "a", To: "b"}}
	out := rename.Apply(entries, rules)
	for i, e := range out {
		if _, ok := e.Fields["a"]; ok {
			t.Errorf("entry %d: old field 'a' should not exist", i)
		}
	}
}
