package output

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/user/logslice/internal/parser"
)

// Format represents the output format type.
type Format string

const (
	FormatJSON  Format = "json"
	FormatText  Format = "text"
	FormatPretty Format = "pretty"
)

// Writer writes log entries to an io.Writer in the specified format.
type Writer struct {
	out    io.Writer
	format Format
}

// New creates a new Writer.
func New(out io.Writer, format Format) *Writer {
	return &Writer{out: out, format: format}
}

// Write writes a single log entry in the configured format.
func (w *Writer) Write(entry parser.Entry) error {
	switch w.format {
	case FormatJSON:
		return w.writeJSON(entry)
	case FormatPretty:
		return w.writePretty(entry)
	default:
		return w.writeText(entry)
	}
}

func (w *Writer) writeJSON(entry parser.Entry) error {
	data, err := json.Marshal(entry.Fields)
	if err != nil {
		return fmt.Errorf("marshal entry: %w", err)
	}
	_, err = fmt.Fprintln(w.out, string(data))
	return err
}

func (w *Writer) writeText(entry parser.Entry) error {
	_, err := fmt.Fprintln(w.out, entry.Raw)
	return err
}

func (w *Writer) writePretty(entry parser.Entry) error {
	var sb strings.Builder
	data, err := json.MarshalIndent(entry.Fields, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal entry: %w", err)
	}
	sb.WriteString(string(data))
	sb.WriteString("\n")
	_, err = fmt.Fprint(w.out, sb.String())
	return err
}
