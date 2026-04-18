package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/user/logslice/internal/grep"
	"github.com/user/logslice/internal/output"
	"github.com/user/logslice/internal/parser"
)

var (
	grepPattern    string
	grepFields     string
	grepIgnoreCase bool
	grepInvert     bool
	grepFormat     string
)

func init() {
	grepCmd := &cobra.Command{
		Use:   "grep [file]",
		Short: "Search log entries by regex pattern",
		Args:  cobra.ExactArgs(1),
		RunE:  runGrep,
	}
	grepCmd.Flags().StringVarP(&grepPattern, "pattern", "p", "", "regex pattern to match (required)")
	grepCmd.Flags().StringVarP(&grepFields, "fields", "f", "", "comma-separated fields to search (default: all)")
	grepCmd.Flags().BoolVarP(&grepIgnoreCase, "ignore-case", "i", false, "case-insensitive matching")
	grepCmd.Flags().BoolVarP(&grepInvert, "invert", "v", false, "invert match")
	grepCmd.Flags().StringVar(&grepFormat, "format", "text", "output format: text, json, pretty")
	_ = grepCmd.MarkFlagRequired("pattern")
	rootCmd.AddCommand(grepCmd)
}

func runGrep(cmd *cobra.Command, args []string) error {
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
	if grepFields != "" {
		for _, fld := range strings.Split(grepFields, ",") {
			if s := strings.TrimSpace(fld); s != "" {
				fields = append(fields, s)
			}
		}
	}

	opts := grep.Options{
		Pattern:     grepPattern,
		Fields:      fields,
		IgnoreCase:  grepIgnoreCase,
		InvertMatch: grepInvert,
	}

	matched, err := grep.Filter(entries, opts)
	if err != nil {
		return fmt.Errorf("grep: %w", err)
	}

	w := output.New(os.Stdout, grepFormat)
	for _, e := range matched {
		w.Write(e)
	}
	return nil
}
