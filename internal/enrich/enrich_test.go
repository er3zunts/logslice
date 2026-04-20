package enrich

import (
	"testing"

	"github.com/user/logslice/internal/parser"
)

func makeEntry(fields map[string]interface{}) parser.Entry {
	return parser.Entry{Fields: fields}
}

func TestParseRules_Valid(t *testing.T) {
	rules, err := ParseRules([]string{"env=prod", "label=${level}-${service}"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(rules) != 2 {
		t.Fatalf("expected 2 rules, got %d", len(rules))
	}
	if rules[0].TargetField != "env" || rules[0].Template != "prod" {
		t.Errorf("unexpected rule[0]: %+v", rules[0])
	}
}

func TestParseRules_Invalid(t *testing.T) {
	_, err := ParseRules([]string{"badspec"})
	if err == nil {
		t.Fatal("expected error for invalid spec")
	}
}

func TestApply_NoRules(t *testing.T) {
	entries := []parser.Entry{makeEntry(map[string]interface{}{"msg": "hello"})}
	out := Apply(entries, nil)
	if len(out) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(out))
	}
}

func TestApply_StaticValue(t *testing.T) {
	entries := []parser.Entry{makeEntry(map[string]interface{}{"msg": "hello"})}
	rules := []Rule{{TargetField: "env", Template: "production"}}
	out := Apply(entries, rules)
	if got := out[0].Fields["env"]; got != "production" {
		t.Errorf("expected 'production', got %v", got)
	}
}

func TestApply_TemplateExpansion(t *testing.T) {
	entries := []parser.Entry{
		makeEntry(map[string]interface{}{"level": "error", "service": "api"}),
	}
	rules := []Rule{{TargetField: "label", Template: "${level}-${service}"}}
	out := Apply(entries, rules)
	if got := out[0].Fields["label"]; got != "error-api" {
		t.Errorf("expected 'error-api', got %v", got)
	}
}

func TestApply_OriginalUnmodified(t *testing.T) {
	original := makeEntry(map[string]interface{}{"msg": "hi"})
	entries := []parser.Entry{original}
	rules := []Rule{{TargetField: "extra", Template: "val"}}
	Apply(entries, rules)
	if _, ok := original.Fields["extra"]; ok {
		t.Error("original entry should not be modified")
	}
}

func TestApply_MultipleRules(t *testing.T) {
	entries := []parser.Entry{
		makeEntry(map[string]interface{}{"host": "web1", "port": "8080"}),
	}
	rules := []Rule{
		{TargetField: "addr", Template: "${host}:${port}"},
		{TargetField: "source", Template: "logslice"},
	}
	out := Apply(entries, rules)
	if got := out[0].Fields["addr"]; got != "web1:8080" {
		t.Errorf("expected 'web1:8080', got %v", got)
	}
	if got := out[0].Fields["source"]; got != "logslice" {
		t.Errorf("expected 'logslice', got %v", got)
	}
}
