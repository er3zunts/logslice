package limit

import (
	"testing"

	"github.com/logslice/logslice/internal/parser"
)

func makeEntry(id string) parser.Entry {
	return parser.Entry{Fields: map[string]interface{}{"id": id}}
}

func entries(ids ...string) []parser.Entry {
	out := make([]parser.Entry, len(ids))
	for i, id := range ids {
		out[i] = makeEntry(id)
	}
	return out
}

func TestApply_NoLimit(t *testing.T) {
	in := entries("a", "b", "c")
	out := Apply(in, Options{Max: 0})
	if len(out) != 3 {
		t.Fatalf("expected 3, got %d", len(out))
	}
}

func TestApply_Limit(t *testing.T) {
	in := entries("a", "b", "c", "d")
	out := Apply(in, Options{Max: 2})
	if len(out) != 2 {
		t.Fatalf("expected 2, got %d", len(out))
	}
	if out[0].Fields["id"] != "a" || out[1].Fields["id"] != "b" {
		t.Errorf("unexpected entries: %v", out)
	}
}

func TestApply_Offset(t *testing.T) {
	in := entries("a", "b", "c", "d")
	out := Apply(in, Options{Offset: 2})
	if len(out) != 2 {
		t.Fatalf("expected 2, got %d", len(out))
	}
	if out[0].Fields["id"] != "c" {
		t.Errorf("expected first entry 'c', got %v", out[0].Fields["id"])
	}
}

func TestApply_OffsetAndLimit(t *testing.T) {
	in := entries("a", "b", "c", "d", "e")
	out := Apply(in, Options{Offset: 1, Max: 2})
	if len(out) != 2 {
		t.Fatalf("expected 2, got %d", len(out))
	}
	if out[0].Fields["id"] != "b" || out[1].Fields["id"] != "c" {
		t.Errorf("unexpected entries: %v", out)
	}
}

func TestApply_OffsetBeyondLength(t *testing.T) {
	in := entries("a", "b")
	out := Apply(in, Options{Offset: 10})
	if len(out) != 0 {
		t.Fatalf("expected 0, got %d", len(out))
	}
}

func TestApply_LimitBeyondLength(t *testing.T) {
	in := entries("a", "b")
	out := Apply(in, Options{Max: 100})
	if len(out) != 2 {
		t.Fatalf("expected 2, got %d", len(out))
	}
}

func TestPage_Basic(t *testing.T) {
	in := entries("a", "b", "c", "d", "e", "f")
	page2 := Page(in, 2, 2)
	if len(page2) != 2 {
		t.Fatalf("expected 2, got %d", len(page2))
	}
	if page2[0].Fields["id"] != "c" || page2[1].Fields["id"] != "d" {
		t.Errorf("unexpected page2 entries: %v", page2)
	}
}

func TestPage_FirstPage(t *testing.T) {
	in := entries("a", "b", "c")
	out := Page(in, 1, 2)
	if len(out) != 2 || out[0].Fields["id"] != "a" {
		t.Errorf("unexpected first page: %v", out)
	}
}

func TestPage_BeyondEnd(t *testing.T) {
	in := entries("a", "b")
	out := Page(in, 5, 3)
	if len(out) != 0 {
		t.Fatalf("expected 0, got %d", len(out))
	}
}
