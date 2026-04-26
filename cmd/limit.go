package cmd

import (
	"fmt"
	"os"

	"github.com/logslice/logslice/internal/limit"
	"github.com/logslice/logslice/internal/output"
	"github.com/logslice/logslice/internal/parser"
	"github.com/spf13/cobra"
)

var (
	limitMax    int
	limitOffset int
	limitPage   int
	limitSize   int
)

func init() {
	limitCmd := &cobra.Command{
		Use:   "limit [file]",
		Short: "Cap the number of log entries returned, with optional offset or pagination",
		Args:  cobra.MaximumNArgs(1),
		RunE:  runLimit,
	}

	limitCmd.Flags().IntVar(&limitMax, "max", 0, "Maximum number of entries to return (0 = no limit)")
	limitCmd.Flags().IntVar(&limitOffset, "offset", 0, "Number of entries to skip before applying --max")
	limitCmd.Flags().IntVar(&limitPage, "page", 0, "1-based page number (requires --page-size; overrides --offset/--max)")
	limitCmd.Flags().IntVar(&limitSize, "page-size", 20, "Number of entries per page")

	rootCmd.AddCommand(limitCmd)
}

func runLimit(cmd *cobra.Command, args []string) error {
	var entries []parser.Entry
	var err error

	if len(args) == 1 {
		f, ferr := os.Open(args[0])
		if ferr != nil {
			return fmt.Errorf("open file: %w", ferr)
		}
		defer f.Close()
		entries, err = parser.ParseJSONLines(f)
	} else {
		entries, err = parser.ParseJSONLines(os.Stdin)
	}
	if err != nil {
		return fmt.Errorf("parse: %w", err)
	}

	var result []parser.Entry
	if limitPage > 0 {
		result = limit.Page(entries, limitPage, limitSize)
	} else {
		result = limit.Apply(entries, limit.Options{
			Max:    limitMax,
			Offset: limitOffset,
		})
	}

	w, err := output.New(outputFormat, os.Stdout)
	if err != nil {
		return fmt.Errorf("output: %w", err)
	}
	for _, e := range result {
		if werr := w.Write(e); werr != nil {
			return werr
		}
	}
	return w.Flush()
}
