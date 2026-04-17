package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/user/logslice/internal/parser"
)

var (
	field  string
	value  string
	input  string
)

var filterCmd = &cobra.Command{
	Use:   "filter",
	Short: "Filter log entries by a field value",
	Example: `  logslice filter --input app.log --field level --value error`,
	RunE: func(cmd *cobra.Command, args []string) error {
		f, err := os.Open(input)
		if err != nil {
			return fmt.Errorf("opening input file: %w", err)
		}
		defer f.Close()

		entries, err := parser.ParseJSONLines(f)
		if err != nil {
			return fmt.Errorf("parsing log file: %w", err)
		}

		for _, entry := range entries {
			if v, ok := entry[field]; ok && fmt.Sprintf("%v", v) == value {
				fmt.Println(entry.String())
			}
		}
		return nil
	},
}

func init() {
	filterCmd.Flags().StringVarP(&field, "field", "f", "", "Field name to filter on (required)")
	filterCmd.Flags().StringVarP(&value, "value", "v", "", "Value to match (required)")
	filterCmd.Flags().StringVarP(&input, "input", "i", "", "Input log file (required)")
	_ = filterCmd.MarkFlagRequired("field")
	_ = filterCmd.MarkFlagRequired("value")
	_ = filterCmd.MarkFlagRequired("input")
}
