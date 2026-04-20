package sortby

import (
	"testing"

	"github.com/user/logslice/internal/parser"
)

func makeEntry(fields map[string]interface{}) parser.Entry {
	return parser.Entry{Fields: fields}
}

func TestApply_Ascending_Numeric(t *testing.T) {
	entries := []parser.Entry{
		makeEntry(map[string]interface{}{"level": "info", "code": 3.0}),
		makeEntry(map[string]interface{}{"level": "warn", "code": 1.0}),
		makeEntry(map[string]interface{}{"level": "error", "code": 2.0}),
	}
	out := Apply(entries, "code", Ascending)
	if out[0].Fields["level"] != "warn" || out[1].Fields["level"] != "error" || out[2].Fields["level"] != "info" {
		t.Errorf("unexpected order: %v", out)
	}
}

func TestApply_Descending_Numeric(t *testing.T) {
	entries := []parser.Entry{
		makeEntry(map[string]interface{}{"code": 1.0}),
		makeEntry(map[string]interface{}{"code": 3.0}),
		makeEntry(map[string]interface{}{"code": 2.0}),
	}
	out := Apply(entries, "code", Descending)
	if out[0].Fields["code"] != 3.0 || out[1].Fields["code"] != 2.0 || out[2].Fields["code"] != 1.0 {
		t.Errorf("unexpected order: %v", out)
	}
}

func TestApply_Ascending_String(t *testing.T) {
	entries := []parser.Entry{
		makeEntry(map[string]interface{}{"msg": "zebra"}),
		makeEntry(map[string]interface{}{"msg": "apple"}),
		makeEntry(map[string]interface{}{"msg": "mango"}),
	}
	out := Apply(entries, "msg", Ascending)
	if out[0].Fields["msg"] != "apple" || out[1].Fields["msg"] != "mango" || out[2].Fields["msg"] != "zebra" {
		t.Errorf("unexpected order: %v", out)
	}
}

func TestApply_MissingFieldSinksToEnd(t *testing.T) {
	entries := []parser.Entry{
		makeEntry(map[string]interface{}{"code": 2.0}),
		makeEntry(map[string]interface{}{}),
		makeEntry(map[string]interface{}{"code": 1.0}),
	}
	out := Apply(entries, "code", Ascending)
	_, hasMissing := out[2].Fields["code"]
	if hasMissing {
		t.Errorf("expected missing-field entry at end, got %v", out[2])
	}
}

func TestApply_DoesNotMutateOriginal(t *testing.T) {
	entries := []parser.Entry{
		makeEntry(map[string]interface{}{"n": 2.0}),
		makeEntry(map[string]interface{}{"n": 1.0}),
	}
	Apply(entries, "n", Ascending)
	if entries[0].Fields["n"] != 2.0 {
		t.Error("original slice was mutated")
	}
}
