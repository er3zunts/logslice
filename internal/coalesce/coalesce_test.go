package coalesce

import (
	"testing"

	"github.com/logslice/logslice/internal/parser"
)

func makeEntry(fields map[string]interface{}) parser.Entry {
	return parser.Entry{Fields: fields}
}

func TestParseRules_Valid(t *testing.T) {
	rules, err := ParseRules([]string{"msg:message,text,body"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(rules) != 1 {
		t.Fatalf("expected 1 rule, got %d", len(rules))
	}
	if rules[0].Dest != "msg" {
		t.Errorf("expected dest=msg, got %s", rules[0].Dest)
	}
	if len(rules[0].Sources) != 3 {
		t.Errorf("expected 3 sources, got %d", len(rules[0].Sources))
	}
}

func TestParseRules_Invalid(t *testing.T) {
	_, err := ParseRules([]string{"nodest"})
	if err == nil {
		t.Fatal("expected error for missing colon, got nil")
	}
}

func TestApply_NoRules(t *testing.T) {
	entries := []parser.Entry{makeEntry(map[string]interface{}{"a": "1"})}
	out := Apply(entries, nil)
	if len(out) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(out))
	}
}

func TestApply_FirstSourceWins(t *testing.T) {
	entries := []parser.Entry{
		makeEntry(map[string]interface{}{"message": "hello", "text": "world"}),
	}
	rules := []Rule{{Dest: "msg", Sources: []string{"message", "text"}}}
	out := Apply(entries, rules)
	if v, ok := out[0].Fields["msg"]; !ok || v != "hello" {
		t.Errorf("expected msg=hello, got %v", v)
	}
}

func TestApply_SkipsEmptyString(t *testing.T) {
	entries := []parser.Entry{
		makeEntry(map[string]interface{}{"message": "", "text": "fallback"}),
	}
	rules := []Rule{{Dest: "msg", Sources: []string{"message", "text"}}}
	out := Apply(entries, rules)
	if v, ok := out[0].Fields["msg"]; !ok || v != "fallback" {
		t.Errorf("expected msg=fallback, got %v", v)
	}
}

func TestApply_SkipsMissingField(t *testing.T) {
	entries := []parser.Entry{
		makeEntry(map[string]interface{}{"text": "only"}),
	}
	rules := []Rule{{Dest: "msg", Sources: []string{"message", "text"}}}
	out := Apply(entries, rules)
	if v, ok := out[0].Fields["msg"]; !ok || v != "only" {
		t.Errorf("expected msg=only, got %v", v)
	}
}

func TestApply_NonePresent(t *testing.T) {
	entries := []parser.Entry{
		makeEntry(map[string]interface{}{"level": "info"}),
	}
	rules := []Rule{{Dest: "msg", Sources: []string{"message", "text"}}}
	out := Apply(entries, rules)
	if _, ok := out[0].Fields["msg"]; ok {
		t.Error("expected msg to be absent when no source found")
	}
}
