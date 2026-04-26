package cmd

import (
	"bufio"
	"fmt"
	"os"

	"github.com/nicholasgasior/logslice/internal/parser"
	"github.com/nicholasgasior/logslice/internal/template"
	"github.com/spf13/cobra"
)

var templateCmd = &cobra.Command{
	Use:   "template",
	Short: "Render log entries using a Go text/template string",
	Long: `Reads JSON log lines from a file (or stdin) and renders each entry
using the provided Go text/template. The entry's fields are available
as template variables (e.g. {{.level}}, {{.message}}).`,
	RunE: runTemplate,
}

var (
	tmplStr  string
	tmplFile string
)

func init() {
	templateCmd.Flags().StringVarP(&tmplStr, "template", "t", "", "Go text/template string (required)")
	templateCmd.Flags().StringVarP(&tmplFile, "file", "f", "", "Input log file (default: stdin)")
	_ = templateCmd.MarkFlagRequired("template")
	rootCmd.AddCommand(templateCmd)
}

func runTemplate(cmd *cobra.Command, args []string) error {
	r, err := template.New(tmplStr)
	if err != nil {
		return fmt.Errorf("invalid template: %w", err)
	}

	var src *os.File
	if tmplFile != "" {
		src, err = os.Open(tmplFile)
		if err != nil {
			return fmt.Errorf("open file: %w", err)
		}
		defer src.Close()
	} else {
		src = os.Stdin
	}

	entries, err := parser.ParseJSONLines(bufio.NewReader(src))
	if err != nil {
		return fmt.Errorf("parse error: %w", err)
	}

	results, renderErrs := r.Apply(entries)
	for _, re := range renderErrs {
		fmt.Fprintf(os.Stderr, "warn: %v\n", re)
	}

	w := bufio.NewWriter(os.Stdout)
	for _, line := range results {
		fmt.Fprintln(w, line)
	}
	return w.Flush()
}
