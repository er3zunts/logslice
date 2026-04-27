package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/user/logslice/internal/parser"
	"github.com/user/logslice/internal/summarize"
)

var (
	summarizeGroupField string
	summarizeValueField string
)

func init() {
	cmd := &cobra.Command{
		Use:   "summarize [file]",
		Short: "Aggregate numeric field statistics grouped by a field",
		Long: `Reads JSON log entries and prints min/max/sum/avg/count for a
numeric field, optionally grouped by another field.

Example:
  logslice summarize -g service -v latency_ms app.log`,
		Args: cobra.MaximumNArgs(1),
		RunE: runSummarize,
	}
	cmd.Flags().StringVarP(&summarizeGroupField, "group", "g", "", "field to group by (omit for global aggregate)")
	cmd.Flags().StringVarP(&summarizeValueField, "value", "v", "", "numeric field to aggregate (required)")
	_ = cmd.MarkFlagRequired("value")
	rootCmd.AddCommand(cmd)
}

func runSummarize(cmd *cobra.Command, args []string) error {
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

	results := summarize.Apply(entries, summarizeGroupField, summarizeValueField)
	if len(results) == 0 {
		fmt.Fprintln(cmd.OutOrStdout(), "no numeric values found for field:", summarizeValueField)
		return nil
	}

	w := cmd.OutOrStdout()
	if summarizeGroupField != "" {
		fmt.Fprintf(w, "%-20s  %8s  %12s  %12s  %12s  %12s\n",
			summarizeGroupField, "count", "sum", "min", "max", "avg")
		fmt.Fprintf(w, "%-20s  %8s  %12s  %12s  %12s  %12s\n",
			"--------------------", "--------", "------------", "------------", "------------", "------------")
	} else {
		fmt.Fprintf(w, "%-20s  %8s  %12s  %12s  %12s  %12s\n",
			"(global)", "count", "sum", "min", "max", "avg")
		fmt.Fprintf(w, "%-20s  %8s  %12s  %12s  %12s  %12s\n",
			"--------------------", "--------", "------------", "------------", "------------", "------------")
	}

	for _, r := range results {
		fmt.Fprintf(w, "%-20s  %8d  %12.4f  %12.4f  %12.4f  %12.4f\n",
			r.GroupKey, r.Stats.Count, r.Stats.Sum, r.Stats.Min, r.Stats.Max, r.Stats.Avg())
	}
	return nil
}
