package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/user/logslice/internal/output"
	"github.com/user/logslice/internal/parser"
	"github.com/user/logslice/internal/tail"
	"github.com/spf13/cobra"
)

var tailLines int
var tailFormat string

var tailCmd = &cobra.Command{
	Use:   "tail [file]",
	Short: "Show the last N log entries from a file",
	Args:  cobra.ExactArgs(1),
	RunE:  runTail,
}

func init() {
	rootCmd.AddCommand(tailCmd)
	tailCmd.Flags().IntVarP(&tailLines, "lines", "n", 10, "Number of lines to show")
	tailCmd.Flags().StringVarP(&tailFormat, "format", "f", "text", "Output format: text, json, pretty")
}

func runTail(cmd *cobra.Command, args []string) error {
	if tailLines <= 0 {
		return fmt.Errorf("lines must be a positive integer, got %d", tailLines)
	}

	path := args[0]
	rawLines, err := tail.Lines(path, tailLines)
	if err != nil {
		return fmt.Errorf("tail: %w", err)
	}

	entries, parseErrs := parser.ParseJSONLines(strings.NewReader(strings.Join(rawLines, "\n")))
	for _, e := range parseErrs {
		fmt.Fprintf(os.Stderr, "warn: %v\n", e)
	}

	w, err := output.New(tailFormat, os.Stdout)
	if err != nil {
		return err
	}
	for _, entry := range entries {
		if err := w.Write(entry); err != nil {
			return err
		}
	}
	return w.Flush()
}
