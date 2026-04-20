package cmd

import (
	"bufio"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/user/logslice/internal/enrich"
	"github.com/user/logslice/internal/output"
	"github.com/user/logslice/internal/parser"
)

var enrichRules []string
var enrichFormat string

func init() {
	enrichCmd := &cobra.Command{
		Use:   "enrich [file]",
		Short: "Add derived fields to log entries using templates",
		Long: `Add new fields to each log entry by expanding templates.

Templates use ${field} syntax to reference existing field values.

Example:
  logslice enrich --rule 'label=${level}-${service}' app.log`,
		Args:    cobra.MaximumNArgs(1),
		RunE:    runEnrich,
	}

	enrichCmd.Flags().StringArrayVar(&enrichRules, "rule", nil, "enrichment rule as field=template (repeatable)")
	enrichCmd.Flags().StringVar(&enrichFormat, "format", "json", "output format: json|text|pretty")
	_ = enrichCmd.MarkFlagRequired("rule")

	rootCmd.AddCommand(enrichCmd)
}

func runEnrich(cmd *cobra.Command, args []string) error {
	rules, err := enrich.ParseRules(enrichRules)
	if err != nil {
		return fmt.Errorf("enrich: %w", err)
	}

	var entries []parser.Entry
	if len(args) == 1 {
		f, err := os.Open(args[0])
		if err != nil {
			return fmt.Errorf("enrich: %w", err)
		}
		defer f.Close()
		entries, err = parser.ParseJSONLines(f)
	} else {
		entries, err = parser.ParseJSONLines(bufio.NewReader(os.Stdin))
	}
	if err != nil {
		return fmt.Errorf("enrich: parse error: %w", err)
	}

	result := enrich.Apply(entries, rules)

	w, err := output.New(os.Stdout, enrichFormat)
	if err != nil {
		return fmt.Errorf("enrich: %w", err)
	}
	for _, e := range result {
		if err := w.Write(e); err != nil {
			return fmt.Errorf("enrich: write error: %w", err)
		}
	}
	return w.Flush()
}
