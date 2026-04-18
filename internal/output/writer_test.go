package output_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/user/logslice/internal/output"
	"github.com/user/logslice/internal/parser"
)

func makeEntry(raw string, fields map[string]interface{}) parser.Entry {
	return parser.Entry{Raw: raw, Fields: fields}
}

func TestWriter_Text(t *testing.T) {
	var buf bytes.Buffer
	w := output.New(&buf, output.FormatText)
	entry := makeEntry(`{"level":"info","msg":"hello"}`, map[string]interface{}{"level": "info", "msg": "hello"})
	if err := w.Write(entry); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got := strings.TrimSpace(buf.String())
	if got != entry.Raw {
		t.Errorf("expected %q, got %q", entry.Raw, got)
	}
}

func TestWriter_JSON(t *testing.T) {
	var buf bytes.Buffer
	w := output.New(&buf, output.FormatJSON)
	entry := makeEntry(`{"level":"info"}`, map[string]interface{}{"level": "info"})
	if err := w.Write(entry); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got := strings.TrimSpace(buf.String())
	if !strings.Contains(got, "level") {
		t.Errorf("expected JSON output to contain 'level', got %q", got)
	}
}

func TestWriter_Pretty(t *testing.T) {
	var buf bytes.Buffer
	w := output.New(&buf, output.FormatPretty)
	entry := makeEntry(`{"level":"warn"}`, map[string]interface{}{"level": "warn"})
	if err := w.Write(entry); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got := buf.String()
	if !strings.Contains(got, "\n") {
		t.Errorf("expected pretty output to contain newlines")
	}
	if !strings.Contains(got, "warn") {
		t.Errorf("expected pretty output to contain 'warn'")
	}
}
