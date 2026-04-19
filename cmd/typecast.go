package cmd

import (
	"bufio"
	"fmt"
	"os"

	"github.com/logslice/logslice/internal/output"
	"github.com/logslice/logslice/internal/parser"
	"github.com/logslice/logslice/internal/typecast"
	"github.com/spf13/cobra"
)

var typecastFields []string
var typecastFormat string

func init() {
	cmd := &cobra.Command{
		Use:   "typecast [file]",
		Short: "Cast field values to specified types",
		Args:  cobra.MaximumNArgs(1),
		RunE:  runTypecast,
	}
	cmd.Flags().StringArrayVarP(&typecastFields, "field", "f", nil, "field:type pairs (e.g. count:int)")
	cmd.Flags().StringVarP(&typecastFormat, "output", "o", "json", "output format: json, text, pretty")
	rootCmd.AddCommand(cmd)
}

func runTypecast(cmd *cobra.Command, args []string) error {
	var r *os.File
	if len(args) == 1 {
		f, err := os.Open(args[0])
		if err != nil {
			return fmt.Errorf("open file: %w", err)
		}
		defer f.Close()
		r = f
	} else {
		r = os.Stdin
	}

	entries, err := parser.ParseJSONLines(bufio.NewReader(r))
	if err != nil {
		return fmt.Errorf("parse: %w", err)
	}

	rules, err := typecast.ParseRules(typecastFields)
	if err != nil {
		return fmt.Errorf("parse rules: %w", err)
	}

	result := typecast.Apply(entries, rules)

	w, err := output.New(typecastFormat, os.Stdout)
	if err != nil {
		return err
	}
	for _, e := range result {
		if err := w.Write(e); err != nil {
			return err
		}
	}
	return nil
}
