package timefmt

import (
	"testing"

	"github.com/nicholasgasior/logslice/internal/parser"
)

func makeEntry(fields map[string]interface{}) parser.Entry {
	return parser.Entry{Fields: fields}
}

func TestParseRules_Valid(t *testing.T) {
	rules, err := ParseRules([]string{"ts:2006-01-02:01/02/2006"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(rules) != 1 || rules[0].Field != "ts" {
		t.Fatalf("unexpected rules: %+v", rules)
	}
}

func TestParseRules_Invalid(t *testing.T) {
	_, err := ParseRules([]string{"badformat"})
	if err == nil {
		t.Fatal("expected error for invalid rule")
	}
}

func TestApply_NoRules(t *testing.T) {
	entries := []parser.Entry{makeEntry(map[string]interface{}{"ts": "2024-01-15"})}
	out := Apply(entries, nil)
	if len(out) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(out))
	}
	if out[0].Fields["ts"] != "2024-01-15" {
		t.Errorf("field should be unchanged")
	}
}

func TestApply_ReformatTimestamp(t *testing.T) {
	entries := []parser.Entry{makeEntry(map[string]interface{}{"ts": "2024-01-15T10:30:00Z"})}
	rules := []Rule{{Field: "ts", InputFmt: "2006-01-02T15:04:05Z", OutputFmt: "02 Jan 2006 15:04"}}
	out := Apply(entries, rules)
	want := "15 Jan 2024 10:30"
	if out[0].Fields["ts"] != want {
		t.Errorf("got %q, want %q", out[0].Fields["ts"], want)
	}
}

func TestApply_MissingField(t *testing.T) {
	entries := []parser.Entry{makeEntry(map[string]interface{}{"level": "info"})}
	rules := []Rule{{Field: "ts", InputFmt: "2006-01-02", OutputFmt: "01/02/2006"}}
	out := Apply(entries, rules)
	if _, ok := out[0].Fields["ts"]; ok {
		t.Error("field should not have been added")
	}
}

func TestApply_UnparseableValue(t *testing.T) {
	entries := []parser.Entry{makeEntry(map[string]interface{}{"ts": "not-a-date"})}
	rules := []Rule{{Field: "ts", InputFmt: "2006-01-02", OutputFmt: "01/02/2006"}}
	out := Apply(entries, rules)
	if out[0].Fields["ts"] != "not-a-date" {
		t.Errorf("unparseable value should be unchanged")
	}
}

func TestApply_NonStringField(t *testing.T) {
	entries := []parser.Entry{makeEntry(map[string]interface{}{"ts": 12345})}
	rules := []Rule{{Field: "ts", InputFmt: "2006-01-02", OutputFmt: "01/02/2006"}}
	out := Apply(entries, rules)
	if out[0].Fields["ts"] != 12345 {
		t.Errorf("non-string value should be unchanged")
	}
}

func TestApply_MultipleRules(t *testing.T) {
	entries := []parser.Entry{makeEntry(map[string]interface{}{
		"start": "2024-03-01T00:00:00Z",
		"end":   "2024-03-02T00:00:00Z",
	})}
	rules := []Rule{
		{Field: "start", InputFmt: "2006-01-02T15:04:05Z", OutputFmt: "2006/01/02"},
		{Field: "end", InputFmt: "2006-01-02T15:04:05Z", OutputFmt: "2006/01/02"},
	}
	out := Apply(entries, rules)
	if out[0].Fields["start"] != "2024/03/01" {
		t.Errorf("start: got %v", out[0].Fields["start"])
	}
	if out[0].Fields["end"] != "2024/03/02" {
		t.Errorf("end: got %v", out[0].Fields["end"])
	}
}
