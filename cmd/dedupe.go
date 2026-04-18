package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/user/logslice/internal/dedupe"
	"github.com/user/logslice/internal/output"
	"github.com/user/logslice/internal/parser"
)

var (
	dedupeFields    string
	dedupeKeepFirst bool
	dedupeFormat    string
)

func init() {
	dedupeCmd := &cobra.Command{
		Use:   "dedupe [file]",
		Short: "Remove duplicate log entries by field values",
		Args:  cobra.ExactArgs(1),
		RunE:  runDedupe,
	}
	dedupeCmd.Flags().StringVar(&dedupeFields, "fields", "", "comma-separated fields to use as dedup key (default: full line)")
	dedupeCmd.Flags().BoolVar(&dedupeKeepFirst, "keep-first", true, "keep first occurrence (false = keep last)")
	dedupeCmd.Flags().StringVar(&dedupeFormat, "format", "text", "output format: text, json, pretty")
	rootCmd.AddCommand(dedupeCmd)
}

func runDedupe(cmd *cobra.Command, args []string) error {
	f, err := os.Open(args[0])
	if err != nil {
		return fmt.Errorf("open: %w", err)
	}
	defer f.Close()

	entries, err := parser.ParseJSONLines(f)
	if err != nil {
		return fmt.Errorf("parse: %w", err)
	}

	var fields []string
	if dedupeFields != "" {
		for _, fld := range strings.Split(dedupeFields, ",") {
			fld = strings.TrimSpace(fld)
			if fld != "" {
				fields = append(fields, fld)
			}
		}
	}

	result := dedupe.Apply(entries, dedupe.Options{
		Fields:    fields,
		KeepFirst: dedupeKeepFirst,
	})

	w := output.New(os.Stdout, dedupeFormat)
	for _, e := range result {
		if err := w.Write(e); err != nil {
			return err
		}
	}
	return w.Flush()
}
