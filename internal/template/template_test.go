package template

import (
	"testing"

	"github.com/nicholasgasior/logslice/internal/parser"
)

func makeEntry(fields map[string]interface{}) parser.Entry {
	return parser.Entry{Fields: fields}
}

func TestNew_ValidTemplate(t *testing.T) {
	_, err := New(`{{.level}} {{.message}}`)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestNew_InvalidTemplate(t *testing.T) {
	_, err := New(`{{.level`)
	if err == nil {
		t.Fatal("expected error for invalid template, got nil")
	}
}

func TestRender_BasicFields(t *testing.T) {
	r, _ := New(`[{{.level}}] {{.message}}`)
	e := makeEntry(map[string]interface{}{"level": "info", "message": "hello"})
	out, err := r.Render(e)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out != "[info] hello" {
		t.Errorf("expected '[info] hello', got %q", out)
	}
}

func TestRender_MissingField(t *testing.T) {
	r, _ := New(`{{.level}} {{.missing}}`)
	e := makeEntry(map[string]interface{}{"level": "warn"})
	out, err := r.Render(e)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// missing keys render as "<no value>" by default in text/template
	if out == "" {
		t.Error("expected non-empty output")
	}
}

func TestRender_DefaultFunc(t *testing.T) {
	r, _ := New(`{{default "n/a" .level}}`)
	e := makeEntry(map[string]interface{}{})
	out, err := r.Render(e)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out != "n/a" {
		t.Errorf("expected 'n/a', got %q", out)
	}
}

func TestApply_MultipleEntries(t *testing.T) {
	r, _ := New(`{{.level}}: {{.msg}}`)
	entries := []parser.Entry{
		makeEntry(map[string]interface{}{"level": "info", "msg": "start"}),
		makeEntry(map[string]interface{}{"level": "error", "msg": "fail"}),
	}
	results, errs := r.Apply(entries)
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
	if results[0] != "info: start" {
		t.Errorf("expected 'info: start', got %q", results[0])
	}
	if results[1] != "error: fail" {
		t.Errorf("expected 'error: fail', got %q", results[1])
	}
}
