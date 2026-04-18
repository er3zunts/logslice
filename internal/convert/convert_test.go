package convert

import (
	"bytes"
	"strings"
	"testing"

	"github.com/user/logslice/internal/parser"
)

func makeEntry(fields map[string]interface{}) parser.Entry {
	return parser.Entry{Fields: fields}
}

func TestToLogfmt_Basic(t *testing.T) {
	entries := []parser.Entry{
		makeEntry(map[string]interface{}{"level": "info", "msg": "hello"}),
	}
	var buf bytes.Buffer
	if err := ToLogfmt(entries, &buf); err != nil {
		t.Fatal(err)
	}
	out := buf.String()
	if !strings.Contains(out, "level=info") {
		t.Errorf("expected level=info in %q", out)
	}
	if !strings.Contains(out, "msg=hello") {
		t.Errorf("expected msg=hello in %q", out)
	}
}

func TestToLogfmt_QuotedValue(t *testing.T) {
	entries := []parser.Entry{
		makeEntry(map[string]interface{}{"msg": "hello world"}),
	}
	var buf bytes.Buffer
	if err := ToLogfmt(entries, &buf); err != nil {
		t.Fatal(err)
	}
	out := buf.String()
	if !strings.Contains(out, `msg=`) {
		t.Errorf("expected quoted msg in %q", out)
	}
}

func TestToCSV_Basic(t *testing.T) {
	entries := []parser.Entry{
		makeEntry(map[string]interface{}{"level": "info", "msg": "start"}),
		makeEntry(map[string]interface{}{"level": "error", "msg": "fail"}),
	}
	var buf bytes.Buffer
	if err := ToCSV(entries, &buf); err != nil {
		t.Fatal(err)
	}
	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(lines) != 3 {
		t.Fatalf("expected 3 lines (header+2), got %d", len(lines))
	}
	if lines[0] != "level,msg" {
		t.Errorf("unexpected header: %q", lines[0])
	}
}

func TestToCSV_Empty(t *testing.T) {
	var buf bytes.Buffer
	if err := ToCSV(nil, &buf); err != nil {
		t.Fatal(err)
	}
	if buf.Len() != 0 {
		t.Errorf("expected empty output for no entries")
	}
}
