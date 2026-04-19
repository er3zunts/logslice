package typecast

import (
	"testing"

	"github.com/logslice/logslice/internal/parser"
)

func makeEntry(fields map[string]interface{}) parser.Entry {
	return parser.Entry{Fields: fields}
}

func TestApply_NoRules(t *testing.T) {
	entries := []parser.Entry{makeEntry(map[string]interface{}{"a": "1"})}
	out := Apply(entries, nil)
	if len(out) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(out))
	}
}

func TestApply_CastToInt(t *testing.T) {
	entries := []parser.Entry{makeEntry(map[string]interface{}{"count": "42"})}
	out := Apply(entries, []Rule{{Field: "count", Target: "int"}})
	val, ok := out[0].Fields["count"]
	if !ok {
		t.Fatal("field missing")
	}
	if _, ok := val.(int64); !ok {
		t.Fatalf("expected int64, got %T", val)
	}
}

func TestApply_CastToFloat(t *testing.T) {
	entries := []parser.Entry{makeEntry(map[string]interface{}{"ratio": "3.14"})}
	out := Apply(entries, []Rule{{Field: "ratio", Target: "float"}})
	if _, ok := out[0].Fields["ratio"].(float64); !ok {
		t.Fatalf("expected float64, got %T", out[0].Fields["ratio"])
	}
}

func TestApply_CastToBool(t *testing.T) {
	entries := []parser.Entry{makeEntry(map[string]interface{}{"ok": "true"})}
	out := Apply(entries, []Rule{{Field: "ok", Target: "bool"}})
	if _, ok := out[0].Fields["ok"].(bool); !ok {
		t.Fatalf("expected bool, got %T", out[0].Fields["ok"])
	}
}

func TestApply_MissingField(t *testing.T) {
	entries := []parser.Entry{makeEntry(map[string]interface{}{"x": "1"})}
	out := Apply(entries, []Rule{{Field: "missing", Target: "int"}})
	if _, ok := out[0].Fields["missing"]; ok {
		t.Fatal("field should not be created")
	}
}

func TestParseRules_Valid(t *testing.T) {
	rules, err := ParseRules([]string{"count:int", "ratio:float"})
	if err != nil {
		t.Fatal(err)
	}
	if len(rules) != 2 || rules[0].Field != "count" || rules[1].Target != "float" {
		t.Fatalf("unexpected rules: %+v", rules)
	}
}

func TestParseRules_Invalid(t *testing.T) {
	_, err := ParseRules([]string{"badspec"})
	if err == nil {
		t.Fatal("expected error")
	}
}
