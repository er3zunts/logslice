package cmd

import (
	"testing"

	"github.com/logslice/logslice/internal/typecast"
)

func TestParseRules_ViaCmd(t *testing.T) {
	rules, err := typecast.ParseRules([]string{"status:int", "latency:float", "ok:bool"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(rules) != 3 {
		t.Fatalf("expected 3 rules, got %d", len(rules))
	}
	expected := []struct{ field, target string }{
		{"status", "int"},
		{"latency", "float"},
		{"ok", "bool"},
	}
	for i, e := range expected {
		if rules[i].Field != e.field || rules[i].Target != e.target {
			t.Errorf("rule %d: got %+v, want %+v", i, rules[i], e)
		}
	}
}

func TestParseRules_EmptySlice(t *testing.T) {
	rules, err := typecast.ParseRules(nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(rules) != 0 {
		t.Fatalf("expected 0 rules, got %d", len(rules))
	}
}
