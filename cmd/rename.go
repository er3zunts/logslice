package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/logslice/logslice/internal/output"
	"github.com/logslice/logslice/internal/parser"
	"github.com/logslice/logslice/internal/rename"
	"github.com/spf13/cobra"
)

var renameMappings []string

func init() {
	renameCmd := &cobra.Command{
		Use:   "rename [file]",
		Short: "Rename fields in log entries",
		Args:  cobra.ExactArgs(1),
		RunE:  runRename,
	}
	renameCmd.Flags().StringArrayVarP(&renameMappings, "map", "m", nil, "Field rename mapping as old=new (repeatable)")
	_ = renameCmd.MarkFlagRequired("map")
	rootCmd.AddCommand(renameCmd)
}

func runRename(cmd *cobra.Command, args []string) error {
	f, err := os.Open(args[0])
	if err != nil {
		return fmt.Errorf("open file: %w", err)
	}
	defer f.Close()

	entries, err := parser.ParseJSONLines(f)
	if err != nil {
		return fmt.Errorf("parse: %w", err)
	}

	rules, err := parseRenameMappings(renameMappings)
	if err != nil {
		return err
	}

	result := rename.Apply(entries, rules)

	w := output.New(os.Stdout, "json")
	for _, e := range result {
		if err := w.Write(e); err != nil {
			return err
		}
	}
	return w.Flush()
}

func parseRenameMappings(mappings []string) ([]rename.Rule, error) {
	rules := make([]rename.Rule, 0, len(mappings))
	for _, m := range mappings {
		parts := strings.SplitN(m, "=", 2)
		if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
			return nil, fmt.Errorf("invalid mapping %q: expected old=new", m)
		}
		rules = append(rules, rename.Rule{From: parts[0], To: parts[1]})
	}
	return rules, nil
}
