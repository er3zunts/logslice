package sample

import (
	"testing"

	"github.com/yourorg/logslice/internal/parser"
)

func makeEntry(msg string) parser.Entry {
	return parser.Entry{Fields: map[string]interface{}{"msg": msg}}
}

func entries(n int) []parser.Entry {
	out := make([]parser.Entry, n)
	for i := range out {
		out[i] = makeEntry("line")
	}
	return out
}

func TestApply_Head(t *testing.T) {
	in := entries(10)
	out := Apply(in, Options{Head: 3})
	if len(out) != 3 {
		t.Fatalf("expected 3, got %d", len(out))
	}
}

func TestApply_Head_MoreThanAvailable(t *testing.T) {
	in := entries(2)
	out := Apply(in, Options{Head: 10})
	if len(out) != 2 {
		t.Fatalf("expected 2, got %d", len(out))
	}
}

func TestApply_Tail(t *testing.T) {
	in := entries(10)
	out := Apply(in, Options{Tail: 4})
	if len(out) != 4 {
		t.Fatalf("expected 4, got %d", len(out))
	}
}

func TestApply_EveryN(t *testing.T) {
	in := entries(10)
	out := Apply(in, Options{Every: 2})
	if len(out) != 5 {
		t.Fatalf("expected 5, got %d", len(out))
	}
}

func TestApply_NoOpts(t *testing.T) {
	in := entries(5)
	out := Apply(in, Options{})
	if len(out) != 5 {
		t.Fatalf("expected 5, got %d", len(out))
	}
}

func TestApply_Rate_Bounds(t *testing.T) {
	in := entries(1000)
	out := Apply(in, Options{Rate: 0.5})
	if len(out) == 0 || len(out) == 1000 {
		t.Fatalf("unexpected count for rate 0.5: %d", len(out))
	}
}
