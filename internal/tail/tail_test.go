package tail

import (
	"strings"
	"testing"
)

func TestReadLastN_AllLines(t *testing.T) {
	input := "line1\nline2\nline3\n"
	lines, err := readLastN(strings.NewReader(input), 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(lines) != 3 {
		t.Fatalf("expected 3 lines, got %d", len(lines))
	}
}

func TestReadLastN_FewerThanN(t *testing.T) {
	input := "a\nb\n"
	lines, err := readLastN(strings.NewReader(input), 10)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(lines) != 2 {
		t.Fatalf("expected 2 lines, got %d", len(lines))
	}
}

func TestReadLastN_ExactN(t *testing.T) {
	input := "x\ny\nz\n"
	lines, err := readLastN(strings.NewReader(input), 2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(lines) != 2 {
		t.Fatalf("expected 2 lines, got %d", len(lines))
	}
	if lines[0] != "y" || lines[1] != "z" {
		t.Fatalf("unexpected lines: %v", lines)
	}
}

func TestReadLastN_Empty(t *testing.T) {
	lines, err := readLastN(strings.NewReader(""), 5)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(lines) != 0 {
		t.Fatalf("expected 0 lines, got %d", len(lines))
	}
}
