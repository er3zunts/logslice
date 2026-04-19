package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/user/logslice/internal/output"
	"github.com/user/logslice/internal/parser"
	"github.com/user/logslice/internal/truncate"
)

var (
	truncMaxLen int
	truncFields string
	truncSuffix string
)

func init() {
	truncCmd := &cobra.Command{
		Use:   "truncate [file]",
		Short: "Truncate long field values in log entries",
		Args:  cobra.ExactArgs(1),
		RunE:  runTruncate,
	}
	truncCmd.Flags().IntVar(&truncMaxLen, "max-len", 80, "maximum character length for field values")
	truncCmd.Flags().StringVar(&truncFields, "fields", "", "comma-separated list of fields to truncate (default: all)")
	truncCmd.Flags().StringVar(&truncSuffix, "suffix", "...", "suffix to append when a value is truncated")
	truncCmd.Flags().StringVar(&outputFormat, "output", "text", "output format: text, json, pretty")
	rootCmd.AddCommand(truncCmd)
}

func runTruncate(cmd *cobra.Command, args []string) error {
	if truncMaxLen <= 0 {
		return fmt.Errorf("--max-len must be a positive integer, got %d", truncMaxLen)
	}

	f, err := os.Open(args[0])
	if err != nil {
		return fmt.Errorf("open file: %w", err)
	}
	defer f.Close()

	entries, err := parser.ParseJSONLines(f)
	if err != nil {
		return fmt.Errorf("parse: %w", err)
	}

	var fields []string
	if truncFields != "" {
		for _, field := range strings.Split(truncFields, ",") {
			if t := strings.TrimSpace(field); t != "" {
				fields = append(fields, t)
			}
		}
	}

	result := truncate.Apply(entries, truncate.Options{
		MaxLen: truncMaxLen,
		Fields: fields,
		Suffix: truncSuffix,
	})

	w := output.New(os.Stdout, outputFormat)
	for _, e := range result {
		if err := w.Write(e); err != nil {
			return err
		}
	}
	return w.Flush()
}
