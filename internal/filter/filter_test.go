package filter

import (
	"testing"
	"time"

	"github.com/user/logslice/internal/parser"
)

func makeEntry(fields map[string]any) parser.Entry {
	return parser.Entry{Fields: fields}
}

func TestApply_NoFilter(t *testing.T) {
	entries := []parser.Entry{
		makeEntry(map[string]any{"msg": "hello"}),
		makeEntry(map[string]any{"msg": "world"}),
	}
	got, err := Apply(entries, Options{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 2 {
		t.Errorf("expected 2 entries, got %d", len(got))
	}
}

func TestApply_FieldFilter(t *testing.T) {
	entries := []parser.Entry{
		makeEntry(map[string]any{"level": "info", "msg": "ok"}),
		makeEntry(map[string]any{"level": "error", "msg": "fail"}),
	}
	got, err := Apply(entries, Options{FieldKey: "level", FieldVal: "error"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 1 || got[0].Fields["msg"] != "fail" {
		t.Errorf("expected 1 error entry, got %+v", got)
	}
}

func TestApply_TimeRange(t *testing.T) {
	from := time.Date(2024, 1, 1, 10, 0, 0, 0, time.UTC)
	to := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)

	entries := []parser.Entry{
		makeEntry(map[string]any{"time": "2024-01-01T09:00:00Z", "msg": "before"}),
		makeEntry(map[string]any{"time": "2024-01-01T11:00:00Z", "msg": "within"}),
		makeEntry(map[string]any{"time": "2024-01-01T13:00:00Z", "msg": "after"}),
	}
	got, err := Apply(entries, Options{TimeFrom: &from, TimeTo: &to})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 1 || got[0].Fields["msg"] != "within" {
		t.Errorf("expected 1 within entry, got %+v", got)
	}
}

func TestApply_SkipsEntriesWithoutTimeField(t *testing.T) {
	from := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	entries := []parser.Entry{
		makeEntry(map[string]any{"msg": "no time field"}),
	}
	got, _ := Apply(entries, Options{TimeFrom: &from})
	if len(got) != 0 {
		t.Errorf("expected 0 entries, got %d", len(got))
	}
}
