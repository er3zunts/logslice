package cmd

import (
	"fmt"
	"os"

	"github.com/logslice/logslice/internal/flatten"
	"github.com/logslice/logslice/internal/output"
	"github.com/logslice/logslice/internal/parser"
	"github.com/spf13/cobra"
)

var (
	flattenSeparator string
	flattenFormat    string
)

func init() {
	flattenCmd := &cobra.Command{
		Use:   "flatten [file]",
		Short: "Flatten nested JSON fields using dot notation",
		Args:  cobra.ExactArgs(1),
		RunE:  runFlatten,
	}
	flattenCmd.Flags().StringVar(&flattenSeparator, "sep", ".", "Separator for nested keys")
	flattenCmd.Flags().StringVar(&flattenFormat, "format", "json", "Output format: json, text, pretty")
	rootCmd.AddCommand(flattenCmd)
}

func runFlatten(cmd *cobra.Command, args []string) error {
	f, err := os.Open(args[0])
	if err != nil {
		return fmt.Errorf("open file: %w", err)
	}
	defer f.Close()

	entries, err := parser.ParseJSONLines(f)
	if err != nil {
		return fmt.Errorf("parse: %w", err)
	}

	flattened := flatten.Apply(entries, flattenSeparator)

	w, err := output.New(os.Stdout, flattenFormat)
	if err != nil {
		return fmt.Errorf("output: %w", err)
	}
	for _, e := range flattened {
		if err := w.Write(e); err != nil {
			return err
		}
	}
	return w.Flush()
}
