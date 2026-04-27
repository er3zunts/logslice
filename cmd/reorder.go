package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/logslice/logslice/internal/output"
	"github.com/logslice/logslice/internal/parser"
	"github.com/logslice/logslice/internal/reorder"
	"github.com/spf13/cobra"
)

var reorderFields string

func init() {
	reorderCmd := &cobra.Command{
		Use:   "reorder [file]",
		Short: "Reorder fields in log entries",
		Long: `Reorder fields in each log entry so that the specified fields appear first,
followed by any remaining fields. Reads from a file or stdin.`,
		RunE: runReorder,
	}
	reorderCmd.Flags().StringVarP(&reorderFields, "fields", "f", "", "Comma-separated list of fields to place first (required)")
	_ = reorderCmd.MarkFlagRequired("fields")
	reorderCmd.Flags().StringVarP(&outputFormat, "output", "o", "json", "Output format: json, text, pretty")
	rootCmd.AddCommand(reorderCmd)
}

func runReorder(cmd *cobra.Command, args []string) error {
	var src *os.File
	if len(args) > 0 {
		f, err := os.Open(args[0])
		if err != nil {
			return fmt.Errorf("opening file: %w", err)
		}
		defer f.Close()
		src = f
	} else {
		src = os.Stdin
	}

	scanner := bufio.NewScanner(src)
	entries, err := parser.ParseJSONLines(scanner)
	if err != nil {
		return fmt.Errorf("parsing input: %w", err)
	}

	fields := splitFieldList(reorderFields)
	result := reorder.Apply(entries, fields)

	w, err := output.New(outputFormat, os.Stdout)
	if err != nil {
		return fmt.Errorf("creating writer: %w", err)
	}
	for _, e := range result {
		if err := w.Write(e); err != nil {
			return fmt.Errorf("writing entry: %w", err)
		}
	}
	return w.Flush()
}

func splitFieldList(s string) []string {
	var out []string
	for _, f := range strings.Split(s, ",") {
		f = strings.TrimSpace(f)
		if f != "" {
			out = append(out, f)
		}
	}
	return out
}
