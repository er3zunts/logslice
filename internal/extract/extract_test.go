package extract

import (
	"testing"

	"github.com/yourorg/logslice/internal/parser"
)

func makeEntry(fields map[string]any) parser.Entry {
	return parser.Entry{Fields: fields}
}

func TestParseRules_Valid(t *testing.T) {
	rules, err := ParseRules([]string{`msg:(?P<word>\w+)`})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(rules) != 1 {
		t.Fatalf("expected 1 rule, got %d", len(rules))
	}
	if rules[0].Field != "msg" {
		t.Errorf("expected field 'msg', got %q", rules[0].Field)
	}
}

func TestParseRules_WithPrefix(t *testing.T) {
	rules, err := ParseRules([]string{`msg:ext_:(?P<word>\w+)`})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if rules[0].Prefix != "ext_" {
		t.Errorf("expected prefix 'ext_', got %q", rules[0].Prefix)
	}
}

func TestParseRules_InvalidRegex(t *testing.T) {
	_, err := ParseRules([]string{`msg:[invalid`})
	if err == nil {
		t.Fatal("expected error for invalid regex")
	}
}

func TestParseRules_NoNamedGroups(t *testing.T) {
	_, err := ParseRules([]string{`msg:(\w+)`})
	if err == nil {
		t.Fatal("expected error for missing named groups")
	}
}

func TestApply_NoRules(t *testing.T) {
	entries := []parser.Entry{makeEntry(map[string]any{"msg": "hello"})}
	out := Apply(entries, nil)
	if len(out) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(out))
	}
}

func TestApply_ExtractsNamedGroups(t *testing.T) {
	rules, _ := ParseRules([]string{`msg:(?P<level>\w+)\s+(?P<code>\d+)`})
	entries := []parser.Entry{makeEntry(map[string]any{"msg": "ERROR 404 not found"})}
	out := Apply(entries, rules)
	if got := out[0].Fields["level"]; got != "ERROR" {
		t.Errorf("expected level=ERROR, got %v", got)
	}
	if got := out[0].Fields["code"]; got != "404" {
		t.Errorf("expected code=404, got %v", got)
	}
}

func TestApply_NoMatch_PassThrough(t *testing.T) {
	rules, _ := ParseRules([]string{`msg:(?P<ip>\d+\.\d+\.\d+\.\d+)`})
	entries := []parser.Entry{makeEntry(map[string]any{"msg": "no ip here"})}
	out := Apply(entries, rules)
	if _, ok := out[0].Fields["ip"]; ok {
		t.Error("expected no 'ip' field on non-matching entry")
	}
}

func TestApply_MissingField_PassThrough(t *testing.T) {
	rules, _ := ParseRules([]string{`msg:(?P<word>\w+)`})
	entries := []parser.Entry{makeEntry(map[string]any{"other": "value"})}
	out := Apply(entries, rules)
	if len(out) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(out))
	}
}

func TestApply_WithPrefix(t *testing.T) {
	rules, _ := ParseRules([]string{`msg:x_:(?P<host>\S+)`})
	entries := []parser.Entry{makeEntry(map[string]any{"msg": "myhost info"})}
	out := Apply(entries, rules)
	if _, ok := out[0].Fields["x_host"]; !ok {
		t.Error("expected field 'x_host' to be present")
	}
}
