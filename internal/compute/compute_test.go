package compute

import (
	"testing"

	"github.com/user/logslice/internal/parser"
)

func makeEntry(fields map[string]interface{}) parser.Entry {
	return parser.Entry{Fields: fields}
}

func TestParseRules_Valid(t *testing.T) {
	rules, err := ParseRules([]string{"diff=end-start", "label=svc+env"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(rules) != 2 {
		t.Fatalf("expected 2 rules, got %d", len(rules))
	}
	if rules[0].Field != "diff" || rules[0].Expr != "end-start" {
		t.Errorf("unexpected rule[0]: %+v", rules[0])
	}
}

func TestParseRules_Invalid(t *testing.T) {
	_, err := ParseRules([]string{"nodivider"})
	if err == nil {
		t.Fatal("expected error for invalid rule")
	}
}

func TestApply_NumericDiff(t *testing.T) {
	entries := []parser.Entry{
		makeEntry(map[string]interface{}{"start": float64(10), "end": float64(35)}),
	}
	rules := []Rule{{Field: "latency", Expr: "end-start"}}
	out := Apply(entries, rules)
	v, ok := out[0].Fields["latency"]
	if !ok {
		t.Fatal("expected latency field")
	}
	if v.(float64) != 25 {
		t.Errorf("expected 25, got %v", v)
	}
}

func TestApply_StringConcat(t *testing.T) {
	entries := []parser.Entry{
		makeEntry(map[string]interface{}{"svc": "auth", "env": "prod"}),
	}
	rules := []Rule{{Field: "label", Expr: "svc+env"}}
	out := Apply(entries, rules)
	if out[0].Fields["label"] != "authprod" {
		t.Errorf("unexpected label: %v", out[0].Fields["label"])
	}
}

func TestApply_NoRules(t *testing.T) {
	entries := []parser.Entry{makeEntry(map[string]interface{}{"x": "y"})}
	out := Apply(entries, nil)
	if len(out) != 1 {
		t.Fatalf("expected 1 entry")
	}
}

func TestApply_MissingField_Skipped(t *testing.T) {
	entries := []parser.Entry{
		makeEntry(map[string]interface{}{"start": float64(5)}),
	}
	rules := []Rule{{Field: "latency", Expr: "end-start"}}
	out := Apply(entries, rules)
	if _, ok := out[0].Fields["latency"]; ok {
		t.Error("expected latency to be absent when operand missing")
	}
}
