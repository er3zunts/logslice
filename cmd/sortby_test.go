package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

func TestRunSort_InvalidFile(t *testing.T) {
	err := runSort(nil, []string{"/nonexistent/path/file.log"})
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestRunSort_ValidFile_Ascending(t *testing.T) {
	content := `{"msg":"third","code":3}
{"msg":"first","code":1}
{"msg":"second","code":2}
`
	tmp := filepath.Join(t.TempDir(), "test.log")
	if err := os.WriteFile(tmp, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	sortField = "code"
	sortDesc = false
	outputFormat = "json"

	// redirect stdout
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	err := runSort(nil, []string{tmp})
	w.Close()
	os.Stdout = old

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var buf bytes.Buffer
	buf.ReadFrom(r)
	out := buf.String()

	// first line should contain "first"
	lines := bytes.Split(bytes.TrimSpace([]byte(out)), []byte("\n"))
	if len(lines) != 3 {
		t.Fatalf("expected 3 lines, got %d", len(lines))
	}
	if !bytes.Contains(lines[0], []byte("first")) {
		t.Errorf("expected first entry to contain 'first', got: %s", lines[0])
	}
}

func TestRunSort_ValidFile_Descending(t *testing.T) {
	content := `{"msg":"third","code":3}
{"msg":"first","code":1}
{"msg":"second","code":2}
`
	tmp := filepath.Join(t.TempDir(), "test.log")
	if err := os.WriteFile(tmp, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	sortField = "code"
	sortDesc = true
	outputFormat = "json"

	// redirect stdout
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	err := runSort(nil, []string{tmp})
	w.Close()
	os.Stdout = old

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var buf bytes.Buffer
	buf.ReadFrom(r)
	out := buf.String()

	// first line should contain "third" when sorted descending
	lines := bytes.Split(bytes.TrimSpace([]byte(out)), []byte("\n"))
	if len(lines) != 3 {
		t.Fatalf("expected 3 lines, got %d", len(lines))
	}
	if !bytes.Contains(lines[0], []byte("third")) {
		t.Errorf("expected first entry to contain 'third', got: %s", lines[0])
	}
}
