package cmd

import (
	"bytes"
	"os"
	"testing"
)

func TestRunTail_InvalidFile(t *testing.T) {
	cmd := tailCmd
	cmd.SetOut(&bytes.Buffer{})
	cmd.SetErr(&bytes.Buffer{})
	err := runTail(cmd, []string{"/nonexistent/path/file.log"})
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}

func TestRunTail_ValidFile(t *testing.T) {
	f, err := os.CreateTemp("", "logslice-tail-*.log")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(f.Name())

	_, _ = f.WriteString(`{"level":"info","msg":"start"}` + "\n")
	_, _ = f.WriteString(`{"level":"warn","msg":"disk low"}` + "\n")
	_, _ = f.WriteString(`{"level":"error","msg":"crash"}` + "\n")
	f.Close()

	tailLines = 2
	tailFormat = "text"

	cmd := tailCmd
	buf := &bytes.Buffer{}
	cmd.SetOut(buf)
	err = runTail(cmd, []string{f.Name()})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
