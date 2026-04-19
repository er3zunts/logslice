package validate

import (
	"regexp"
	"testing"

	"github.com/user/logslice/internal/parser"
)

func makeEntry(fields map[string]interface{}) parser.Entry {
	return parser.Entry{Fields: fields}
}

func TestApply_NoRules(t *testing.T) {
	entries := []parser.Entry{makeEntry(map[string]interface{}{"level": "info"})}
	results := Apply(entries, nil)
	if len(results) != 1 || !results[0].Valid() {
		t.Fatal("expected valid result with no rules")
	}
}

func TestApply_RequiredFieldPresent(t *testing.T) {
	entries := []parser.Entry{makeEntry(map[string]interface{}{"level": "info"})}
	rules := []Rule{{Field: "level", Required: true}}
	results := Apply(entries, rules)
	if !results[0].Valid() {
		t.Fatal("expected valid")
	}
}

func TestApply_RequiredFieldMissing(t *testing.T) {
	entries := []parser.Entry{makeEntry(map[string]interface{}{"msg": "hello"})}
	rules := []Rule{{Field: "level", Required: true}}
	results := Apply(entries, rules)
	if results[0].Valid() {
		t.Fatal("expected invalid")
	}
	if len(results[0].Errors) != 1 {
		t.Fatalf("expected 1 error, got %d", len(results[0].Errors))
	}
}

func TestApply_PatternMatch(t *testing.T) {
	entries := []parser.Entry{
		makeEntry(map[string]interface{}{"level": "info"}),
		makeEntry(map[string]interface{}{"level": "bad-level"}),
	}
	pat := regexp.MustCompile(`^(info|warn|error)$`)
	rules := []Rule{{Field: "level", Pattern: pat}}
	results := Apply(entries, rules)
	if !results[0].Valid() {
		t.Error("expected first entry valid")
	}
	if results[1].Valid() {
		t.Error("expected second entry invalid")
	}
}

func TestApply_TypeCheck(t *testing.T) {
	entries := []parser.Entry{
		makeEntry(map[string]interface{}{"count": float64(3)}),
		makeEntry(map[string]interface{}{"count": "not-a-number"}),
	}
	rules := []Rule{{Field: "count", TypeName: "number"}}
	results := Apply(entries, rules)
	if !results[0].Valid() {
		t.Error("expected first entry valid")
	}
	if results[1].Valid() {
		t.Error("expected second entry invalid")
	}
}

func TestResult_Valid(t *testing.T) {
	r := Result{Errors: []string{}}
	if !r.Valid() {
		t.Fatal("expected valid")
	}
	r.Errors = append(r.Errors, "some error")
	if r.Valid() {
		t.Fatal("expected invalid")
	}
}
