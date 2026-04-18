package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/user/logslice/internal/output"
	"github.com/user/logslice/internal/parser"
)

var (
	outputFormat string
	outputFile   string
)

func init() {
	outputCmd := &cobra.Command{
		Use:   "output [file]",
		Short: "Re-format a log file and write to stdout or a file",
		Args:  cobra.ExactArgs(1),
		RunE:  runOutput,
	}
	outputCmd.Flags().StringVarP(&outputFormat, "format", "f", "text", "Output format: text, json, pretty")
	outputCmd.Flags().StringVarP(&outputFile, "out", "o", "", "Output file (default: stdout)")
	rootCmd.AddCommand(outputCmd)
}

func runOutput(cmd *cobra.Command, args []string) error {
	f, err := os.Open(args[0])
	if err != nil {
		return fmt.Errorf("open file: %w", err)
	}
	defer f.Close()

	entries, parseErrs := parser.ParseJSONLines(f)
	for _, e := range parseErrs {
		fmt.Fprintf(os.Stderr, "warning: %v\n", e)
	}

	dest := cmd.OutOrStdout()
	if outputFile != "" {
		out, err := os.Create(outputFile)
		if err != nil {
			return fmt.Errorf("create output file: %w", err)
		}
		defer out.Close()
		dest = out
	}

	w := output.New(dest, output.Format(outputFormat))
	for _, entry := range entries {
		if err := w.Write(entry); err != nil {
			return fmt.Errorf("write entry: %w", err)
		}
	}
	return nil
}
