package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/user/logslice/internal/convert"
	"github.com/user/logslice/internal/parser"
)

var convertFormat string

func init() {
	convertCmd := &cobra.Command{
		Use:   "convert [file]",
		Short: "Convert log entries to a different format (csv, logfmt)",
		Args:  cobra.ExactArgs(1),
		RunE:  runConvert,
	}
	convertCmd.Flags().StringVarP(&convertFormat, "format", "f", "logfmt", "output format: csv or logfmt")
	rootCmd.AddCommand(convertCmd)
}

func runConvert(cmd *cobra.Command, args []string) error {
	f, err := os.Open(args[0])
	if err != nil {
		return fmt.Errorf("open file: %w", err)
	}
	defer f.Close()

	entries, err := parser.ParseJSONLines(f)
	if err != nil {
		return fmt.Errorf("parse: %w", err)
	}

	switch convert.Format(convertFormat) {
	case convert.FormatCSV:
		return convert.ToCSV(entries, os.Stdout)
	case convert.FormatLogfmt:
		return convert.ToLogfmt(entries, os.Stdout)
	default:
		return fmt.Errorf("unsupported format %q: use csv or logfmt", convertFormat)
	}
}
