package merge

import (
	"testing"

	"github.com/user/logslice/internal/parser"
)

func makeEntry(ts, msg string) parser.Entry {
	fields := map[string]interface{}{"msg": msg}
	if ts != "" {
		fields[TimeKey] = ts
	}
	return parser.Entry{Fields: fields, Raw: msg}
}

func TestByTime_SingleStream(t *testing.T) {
	stream := []parser.Entry{
		makeEntry("2024-01-01T10:00:00Z", "a"),
		makeEntry("2024-01-01T09:00:00Z", "b"),
	}
	result := ByTime(stream)
	if result[0].Fields["msg"] != "b" || result[1].Fields["msg"] != "a" {
		t.Errorf("expected sorted by time, got %v %v", result[0].Fields["msg"], result[1].Fields["msg"])
	}
}

func TestByTime_MultipleStreams(t *testing.T) {
	s1 := []parser.Entry{
		makeEntry("2024-01-01T08:00:00Z", "s1-early"),
		makeEntry("2024-01-01T12:00:00Z", "s1-late"),
	}
	s2 := []parser.Entry{
		makeEntry("2024-01-01T10:00:00Z", "s2-mid"),
	}
	result := ByTime(s1, s2)
	if len(result) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(result))
	}
	expected := []string{"s1-early", "s2-mid", "s1-late"}
	for i, e := range result {
		if e.Fields["msg"] != expected[i] {
			t.Errorf("pos %d: expected %s, got %v", i, expected[i], e.Fields["msg"])
		}
	}
}

func TestByTime_NoTimestamp(t *testing.T) {
	stream := []parser.Entry{
		makeEntry("", "no-ts"),
		makeEntry("2024-01-01T10:00:00Z", "with-ts"),
	}
	result := ByTime(stream)
	if result[0].Fields["msg"] != "with-ts" {
		t.Errorf("expected timestamped entry first, got %v", result[0].Fields["msg"])
	}
}

func TestByTime_Empty(t *testing.T) {
	result := ByTime([]parser.Entry{}, []parser.Entry{})
	if len(result) != 0 {
		t.Errorf("expected empty result, got %d", len(result))
	}
}
