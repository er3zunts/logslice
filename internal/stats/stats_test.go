package stats

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/user/logslice/internal/parser"
)

func makeEntry(fields map[string]interface{}) parser.Entry {
	return parser.Entry{Time: time.Now(), Fields: fields}
}

func TestCompute_Total(t *testing.T) {
	entries := []parser.Entry{makeEntry(nil), makeEntry(nil)}
	s := Compute(entries, nil)
	if s.Total != 2 {
		t.Errorf("expected 2, got %d", s.Total)
	}
}

func TestCompute_FieldCounts(t *testing.T) {
	entries := []parser.Entry{
		makeEntry(map[string]interface{}{"level": "info"}),
		makeEntry(map[string]interface{}{"level": "error"}),
		makeEntry(map[string]interface{}{"level": "info"}),
	}
	s := Compute(entries, []string{"level"})
	if s.FieldCounts["level"]["info"] != 2 {
		t.Errorf("expected info=2, got %d", s.FieldCounts["level"]["info"])
	}
	if s.FieldCounts["level"]["error"] != 1 {
		t.Errorf("expected error=1, got %d", s.FieldCounts["level"]["error"])
	}
}

func TestCompute_MissingField(t *testing.T) {
	entries := []parser.Entry{
		makeEntry(map[string]interface{}{"level": "info"}),
		makeEntry(map[string]interface{}{}),
	}
	s := Compute(entries, []string{"level"})
	if s.FieldCounts["level"]["info"] != 1 {
		t.Errorf("expected 1, got %d", s.FieldCounts["level"]["info"])
	}
}

func TestSummary_Print(t *testing.T) {
	entries := []parser.Entry{
		makeEntry(map[string]interface{}{"level": "info"}),
	}
	s := Compute(entries, []string{"level"})
	var buf bytes.Buffer
	s.Print(&buf)
	out := buf.String()
	if !strings.Contains(out, "Total entries: 1") {
		t.Errorf("missing total in output: %s", out)
	}
	if !strings.Contains(out, "info") {
		t.Errorf("missing field value in output: %s", out)
	}
}
