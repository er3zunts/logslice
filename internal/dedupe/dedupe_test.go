package dedupe

import (
	"testing"

	"github.com/user/logslice/internal/parser"
)

func makeEntry(raw string, fields map[string]interface{}) parser.Entry {
	return parser.Entry{Raw: raw, Fields: fields}
}

func TestApply_NoDuplicates(t *testing.T) {
	entries := []parser.Entry{
		makeEntry(`{"msg":"a"}`, map[string]interface{}{"msg": "a"}),
		makeEntry(`{"msg":"b"}`, map[string]interface{}{"msg": "b"}),
	}
	out := Apply(entries, Options{})
	if len(out) != 2 {
		t.Fatalf("expected 2, got %d", len(out))
	}
}

func TestApply_KeepFirst(t *testing.T) {
	entries := []parser.Entry{
		makeEntry(`{"msg":"a","v":1}`, map[string]interface{}{"msg": "a", "v": float64(1)}),
		makeEntry(`{"msg":"a","v":2}`, map[string]interface{}{"msg": "a", "v": float64(2)}),
	}
	out := Apply(entries, Options{Fields: []string{"msg"}, KeepFirst: true})
	if len(out) != 1 {
		t.Fatalf("expected 1, got %d", len(out))
	}
	if out[0].Fields["v"] != float64(1) {
		t.Errorf("expected first entry kept")
	}
}

func TestApply_KeepLast(t *testing.T) {
	entries := []parser.Entry{
		makeEntry(`{"msg":"a","v":1}`, map[string]interface{}{"msg": "a", "v": float64(1)}),
		makeEntry(`{"msg":"a","v":2}`, map[string]interface{}{"msg": "a", "v": float64(2)}),
	}
	out := Apply(entries, Options{Fields: []string{"msg"}, KeepFirst: false})
	if len(out) != 1 {
		t.Fatalf("expected 1, got %d", len(out))
	}
	if out[0].Fields["v"] != float64(2) {
		t.Errorf("expected last entry kept")
	}
}

func TestApply_MultipleFields(t *testing.T) {
	entries := []parser.Entry{
		makeEntry(`1`, map[string]interface{}{"svc": "api", "code": "200"}),
		makeEntry(`2`, map[string]interface{}{"svc": "api", "code": "500"}),
		makeEntry(`3`, map[string]interface{}{"svc": "api", "code": "200"}),
	}
	out := Apply(entries, Options{Fields: []string{"svc", "code"}, KeepFirst: true})
	if len(out) != 2 {
		t.Fatalf("expected 2, got %d", len(out))
	}
}
