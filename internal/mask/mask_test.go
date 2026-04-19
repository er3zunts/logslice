package mask

import (
	"regexp"
	"testing"

	"github.com/user/logslice/internal/parser"
)

func makeEntry(fields map[string]interface{}) parser.Entry {
	return parser.Entry{Fields: fields}
}

func TestApply_MaskWholeField(t *testing.T) {
	entries := []parser.Entry{
		makeEntry(map[string]interface{}{"password": "secret123", "user": "alice"}),
	}
	rules := []Rule{{Field: "password"}}
	out := Apply(entries, rules)
	if got := out[0].Fields["password"]; got != "*********" {
		t.Errorf("expected masked password, got %q", got)
	}
	if got := out[0].Fields["user"]; got != "alice" {
		t.Errorf("user should be unchanged, got %q", got)
	}
}

func TestApply_MaskWithPattern(t *testing.T) {
	entries := []parser.Entry{
		makeEntry(map[string]interface{}{"email": "user@example.com"}),
	}
	rules := []Rule{
		{Field: "email", Pattern: regexp.MustCompile(`[^@]+@`), Replace: "***@"},
	}
	out := Apply(entries, rules)
	if got := out[0].Fields["email"]; got != "***@example.com" {
		t.Errorf("unexpected email mask: %q", got)
	}
}

func TestApply_MissingField(t *testing.T) {
	entries := []parser.Entry{
		makeEntry(map[string]interface{}{"msg": "hello"}),
	}
	rules := []Rule{{Field: "token"}}
	out := Apply(entries, rules)
	if _, ok := out[0].Fields["token"]; ok {
		t.Error("token field should not be added")
	}
}

func TestApply_NonStringUnchanged(t *testing.T) {
	entries := []parser.Entry{
		makeEntry(map[string]interface{}{"count": 42.0}),
	}
	rules := []Rule{{Field: "count"}}
	out := Apply(entries, rules)
	if got := out[0].Fields["count"]; got != 42.0 {
		t.Errorf("non-string field should be unchanged, got %v", got)
	}
}

func TestParseRules_NoPattern(t *testing.T) {
	rules, err := ParseRules([]string{"password", "token"})
	if err != nil {
		t.Fatal(err)
	}
	if len(rules) != 2 || rules[0].Field != "password" || rules[0].Pattern != nil {
		t.Errorf("unexpected rules: %+v", rules)
	}
}

func TestParseRules_WithPattern(t *testing.T) {
	rules, err := ParseRules([]string{`email=[^@]+@`})
	if err != nil {
		t.Fatal(err)
	}
	if rules[0].Pattern == nil {
		t.Error("expected pattern to be set")
	}
}

func TestParseRules_InvalidPattern(t *testing.T) {
	_, err := ParseRules([]string{"field=[invalid"})
	if err == nil {
		t.Error("expected error for invalid regex")
	}
}
