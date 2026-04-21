package cmd

import (
	"fmt"
	"os"

	"github.com/logslice/logslice/internal/groupby"
	"github.com/logslice/logslice/internal/parser"
	"github.com/spf13/cobra"
)

var (
	groupField    string
	groupKeepRows bool
)

func init() {
	groupCmd := &cobra.Command{
		Use:   "groupby [file]",
		Short: "Group log entries by a field value and show counts",
		Args:  cobra.ExactArgs(1),
		RunE:  runGroupBy,
	}
	groupCmd.Flags().StringVarP(&groupField, "field", "f", "level", "Field to group by")
	groupCmd.Flags().BoolVar(&groupKeepRows, "keep-rows", false, "Print individual entries under each group")
	rootCmd.AddCommand(groupCmd)
}

func runGroupBy(cmd *cobra.Command, args []string) error {
	f, err := os.Open(args[0])
	if err != nil {
		return fmt.Errorf("open file: %w", err)
	}
	defer f.Close()

	entries, err := parser.ParseJSONLines(f)
	if err != nil {
		return fmt.Errorf("parse: %w", err)
	}

	groups := groupby.Apply(entries, groupby.Options{
		Field:    groupField,
		KeepRows: groupKeepRows,
	})

	for _, g := range groups {
		fmt.Fprintf(cmd.OutOrStdout(), "%s\t%d\n", g.Key, g.Count)
		if groupKeepRows {
			for _, e := range g.Items {
				fmt.Fprintf(cmd.OutOrStdout(), "  %v\n", e.Fields)
			}
		}
	}
	return nil
}
