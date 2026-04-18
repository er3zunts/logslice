package grep_test

import (
	"testing"

	"github.com/user/logslice/internal/grep"
	"github.com/user/logslice/internal/parser"
)

func makeEntry(fields map[string]string) parser.Entry {
	return parser.Entry{Fields: fields}
}

func TestMatch_AllFields(t *testing.T) {
	e := makeEntry(map[string]string{"msg": "server started", "level": "info"})
	ok, err := grep.Match(e, grep.Options{Pattern: "started"})
	if err != nil || !ok {
		t.Fatalf("expected match, got ok=%v err=%v", ok, err)
	}
}

func TestMatch_SpecificField(t *testing.T) {
	e := makeEntry(map[string]string{"msg": "hello world", "level": "error"})
	ok, err := grep.Match(e, grep.Options{Pattern: "error", Fields: []string{"level"}})
	if err != nil || !ok {
		t.Fatalf("expected match")
	}
	ok2, _ := grep.Match(e, grep.Options{Pattern: "error", Fields: []string{"msg"}})
	if ok2 {
		t.Fatal("should not match msg field")
	}
}

func TestMatch_IgnoreCase(t *testing.T) {
	e := makeEntry(map[string]string{"msg": "Hello World"})
	ok, err := grep.Match(e, grep.Options{Pattern: "hello", IgnoreCase: true})
	if err != nil || !ok {
		t.Fatal("expected case-insensitive match")
	}
}

func TestMatch_InvertMatch(t *testing.T) {
	e := makeEntry(map[string]string{"level": "debug"})
	ok, err := grep.Match(e, grep.Options{Pattern: "error", InvertMatch: true})
	if err != nil || !ok {
		t.Fatal("expected inverted match")
	}
}

func TestMatch_InvalidPattern(t *testing.T) {
	e := makeEntry(map[string]string{"msg": "test"})
	_, err := grep.Match(e, grep.Options{Pattern: "[invalid"})
	if err == nil {
		t.Fatal("expected error for invalid regex")
	}
}

func TestFilter(t *testing.T) {
	entries := []parser.Entry{
		makeEntry(map[string]string{"msg": "started"}),
		makeEntry(map[string]string{"msg": "stopped"}),
		makeEntry(map[string]string{"msg": "restarted"}),
	}
	result, err := grep.Filter(entries, grep.Options{Pattern: "start"})
	if err != nil {
		t.Fatal(err)
	}
	if len(result) != 2 {
		t.Fatalf("expected 2 matches, got %d", len(result))
	}
}
