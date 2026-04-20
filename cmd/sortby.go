package cmd

import (
	"bufio"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/user/logslice/internal/output"
	"github.com/user/logslice/internal/parser"
	"github.com/user/logslice/internal/sortby"
)

var (
	sortField string
	sortDesc  bool
)

func init() {
	sortCmd := &cobra.Command{
		Use:   "sort [file]",
		Short: "Sort log entries by a field value",
		Args:  cobra.MaximumNArgs(1),
		RunE:  runSort,
	}
	sortCmd.Flags().StringVarP(&sortField, "field", "f", "timestamp", "Field name to sort by")
	sortCmd.Flags().BoolVarP(&sortDesc, "desc", "d", false, "Sort in descending order")
	sortCmd.Flags().StringVarP(&outputFormat, "output", "o", "text", "Output format: text, json, pretty")
	rootCmd.AddCommand(sortCmd)
}

func runSort(cmd *cobra.Command, args []string) error {
	var src *os.File
	if len(args) == 1 {
		f, err := os.Open(args[0])
		if err != nil {
			return fmt.Errorf("open file: %w", err)
		}
		defer f.Close()
		src = f
	} else {
		src = os.Stdin
	}

	entries, err := parser.ParseJSONLines(bufio.NewReader(src))
	if err != nil {
		return fmt.Errorf("parse: %w", err)
	}

	order := sortby.Ascending
	if sortDesc {
		order = sortby.Descending
	}

	sorted := sortby.Apply(entries, sortField, order)

	w, err := output.New(outputFormat, os.Stdout)
	if err != nil {
		return err
	}
	for _, e := range sorted {
		if err := w.Write(e); err != nil {
			return err
		}
	}
	return nil
}
