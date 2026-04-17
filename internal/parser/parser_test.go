package parser

import (
	"strings"
	"testing"
)

func TestParseJSONLines_ValidInput(t *testing.T) {
	input := `{"level":"info","msg":"started"}
{"level":"error","msg":"failed","code":500}
`
	entries, err := ParseJSONLines(strings.NewReader(input))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}
	if entries[0]["level"] != "info" {
		t.Errorf("expected level=info, got %v", entries[0]["level"])
	}
	if entries[1]["msg"] != "failed" {
		t.Errorf("expected msg=failed, got %v", entries[1]["msg"])
	}
}

func TestParseJSONLines_SkipsBlankLines(t *testing.T) {
	input := `{"level":"warn"}

{"level":"debug"}
`
	entries, err := ParseJSONLines(strings.NewReader(input))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}
}

func TestParseJSONLines_InvalidJSON(t *testing.T) {
	input := `{"level":"info"}
not-json
`
	_, err := ParseJSONLines(strings.NewReader(input))
	if err == nil {
		t.Fatal("expected error for invalid JSON, got nil")
	}
}

func TestEntryString(t *testing.T) {
	e := Entry{"level": "info", "msg": "hello"}
	s := e.String()
	if s == "" {
		t.Error("expected non-empty string from Entry.String()")
	}
}
