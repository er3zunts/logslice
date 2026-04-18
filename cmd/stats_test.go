package cmd

import (
	"testing"
)

func TestSplitFields_Basic(t *testing.T) {
	result := splitFields("level,service,host")
	if len(result) != 3 {
		t.Fatalf("expected 3 fields, got %d", len(result))
	}
	if result[0] != "level" || result[1] != "service" || result[2] != "host" {
		t.Errorf("unexpected fields: %v", result)
	}
}

func TestSplitFields_Spaces(t *testing.T) {
	result := splitFields(" level , service ")
	if len(result) != 2 {
		t.Fatalf("expected 2, got %d", len(result))
	}
	if result[0] != "level" {
		t.Errorf("expected 'level', got '%s'", result[0])
	}
}

func TestSplitFields_Empty(t *testing.T) {
	result := splitFields("")
	if len(result) != 0 {
		t.Errorf("expected empty slice, got %v", result)
	}
}

func TestSplitFields_Single(t *testing.T) {
	result := splitFields("level")
	if len(result) != 1 || result[0] != "level" {
		t.Errorf("unexpected: %v", result)
	}
}
