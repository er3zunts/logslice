package fieldmap

import (
	"testing"

	"github.com/logslice/logslice/internal/parser"
)

func makeEntry(fields map[string]interface{}) parser.Entry {
	return parser.Entry{Fields: fields}
}

func TestApply_NoRules(t *testing.T) {
	entries := []parser.Entry{makeEntry(map[string]interface{}{"a": 1})}
	out := Apply(entries, nil)
	if len(out) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(out))
	}
}

func TestApply_SelectFields(t *testing.T) {
	e := makeEntry(map[string]interface{}{"level": "info", "msg": "hello", "ts": "now"})
	out := Apply([]parser.Entry{e}, []Rule{{Field: "level"}, {Field: "msg"}})
	if len(out) != 1 {
		t.Fatalf("expected 1 entry")
	}
	if _, ok := out[0].Fields["ts"]; ok {
		t.Error("ts should have been excluded")
	}
	if out[0].Fields["level"] != "info" {
		t.Errorf("expected level=info")
	}
}

func TestApply_AliasField(t *testing.T) {
	e := makeEntry(map[string]interface{}{"level": "warn"})
	out := Apply([]parser.Entry{e}, []Rule{{Field: "level", Alias: "severity"}})
	if _, ok := out[0].Fields["level"]; ok {
		t.Error("original field name should not appear")
	}
	if out[0].Fields["severity"] != "warn" {
		t.Errorf("expected severity=warn")
	}
}

func TestApply_MissingField(t *testing.T) {
	e := makeEntry(map[string]interface{}{"msg": "hi"})
	out := Apply([]parser.Entry{e}, []Rule{{Field: "level"}})
	if _, ok := out[0].Fields["level"]; ok {
		t.Error("missing field should not appear in output")
	}
}

func TestParseRules_Basic(t *testing.T) {
	rules := ParseRules([]string{"level", "msg:message"})
	if len(rules) != 2 {
		t.Fatalf("expected 2 rules")
	}
	if rules[0].Field != "level" || rules[0].Alias != "" {
		t.Errorf("unexpected rule[0]: %+v", rules[0])
	}
	if rules[1].Field != "msg" || rules[1].Alias != "message" {
		t.Errorf("unexpected rule[1]: %+v", rules[1])
	}
}
