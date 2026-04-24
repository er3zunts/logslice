package cmd

import (
	"fmt"
	"os"

	"github.com/nicholasgasior/logslice/internal/output"
	"github.com/nicholasgasior/logslice/internal/parser"
	"github.com/nicholasgasior/logslice/internal/timefmt"
	"github.com/spf13/cobra"
)

var (
	timefmtRules  []string
	timefmtFormat string
)

func init() {
	timefmtCmd := &cobra.Command{
		Use:   "timefmt [file]",
		Short: "Reformat timestamp fields in log entries",
		Long: `Reformat one or more timestamp fields using Go time layouts.

Rules are specified as: field:inputLayout:outputLayout

Example:
  logslice timefmt app.log --rule "ts:2006-01-02T15:04:05Z07:00:02 Jan 2006"
`,
		Args: cobra.ExactArgs(1),
		RunE: runTimefmt,
	}
	timefmtCmd.Flags().StringArrayVar(&timefmtRules, "rule", nil, "reformat rule: field:inputFmt:outputFmt (repeatable)")
	timefmtCmd.Flags().StringVar(&timefmtFormat, "output", "json", "output format: json|text|pretty")
	_ = timefmtCmd.MarkFlagRequired("rule")
	rootCmd.AddCommand(timefmtCmd)
}

func runTimefmt(cmd *cobra.Command, args []string) error {
	f, err := os.Open(args[0])
	if err != nil {
		return fmt.Errorf("timefmt: open file: %w", err)
	}
	defer f.Close()

	entries, err := parser.ParseJSONLines(f)
	if err != nil {
		return fmt.Errorf("timefmt: parse: %w", err)
	}

	rules, err := timefmt.ParseRules(timefmtRules)
	if err != nil {
		return err
	}

	result := timefmt.Apply(entries, rules)

	w, err := output.New(timefmtFormat, os.Stdout)
	if err != nil {
		return fmt.Errorf("timefmt: output: %w", err)
	}
	for _, e := range result {
		if err := w.Write(e); err != nil {
			return err
		}
	}
	return w.Flush()
}
