package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/user/logslice/internal/mask"
	"github.com/user/logslice/internal/output"
	"github.com/user/logslice/internal/parser"
)

var (
	maskFields  []string
	maskFormat  string
)

func init() {
	maskCmd := &cobra.Command{
		Use:   "mask [file]",
		Short: "Mask sensitive field values in log entries",
		Long:  "Mask replaces sensitive field values with asterisks or a regex-based replacement.\nSpecify fields as 'fieldname' or 'fieldname=pattern'.",
		Args:  cobra.ExactArgs(1),
		RunE:  runMask,
	}
	maskCmd.Flags().StringArrayVarP(&maskFields, "field", "f", nil, "field to mask, optionally with pattern: field=pattern (repeatable)")
	maskCmd.Flags().StringVarP(&maskFormat, "output", "o", "json", "output format: json, text, pretty")
	_ = maskCmd.MarkFlagRequired("field")
	rootCmd.AddCommand(maskCmd)
}

func runMask(cmd *cobra.Command, args []string) error {
	f, err := os.Open(args[0])
	if err != nil {
		return fmt.Errorf("open file: %w", err)
	}
	defer f.Close()

	entries, err := parser.ParseJSONLines(f)
	if err != nil {
		return fmt.Errorf("parse: %w", err)
	}

	rules, err := mask.ParseRules(maskFields)
	if err != nil {
		return fmt.Errorf("parse mask rules: %w", err)
	}

	result := mask.Apply(entries, rules)

	w, err := output.New(maskFormat, os.Stdout)
	if err != nil {
		return fmt.Errorf("output: %w", err)
	}
	for _, e := range result {
		if err := w.Write(e); err != nil {
			return err
		}
	}
	return w.Flush()
}
