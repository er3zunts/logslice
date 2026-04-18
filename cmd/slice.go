package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
	"github.com/user/logslice/internal/filter"
	"github.com/user/logslice/internal/parser"
)

var (
	sliceFrom      string
	sliceTo        string
	sliceField     string
	sliceFieldVal  string
	sliceTimeField string
)

var sliceCmd = &cobra.Command{
	Use:   "slice [file]",
	Short: "Slice log entries by time range and/or field value",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		f, err := os.Open(args[0])
		if err != nil {
			return fmt.Errorf("opening file: %w", err)
		}
		defer f.Close()

		entries, err := parser.ParseJSONLines(f)
		if err != nil {
			return fmt.Errorf("parsing: %w", err)
		}

		opts := filter.Options{
			FieldKey:  sliceField,
			FieldVal:  sliceFieldVal,
			TimeField: sliceTimeField,
		}

		if sliceFrom != "" {
			t, err := time.Parse(time.RFC3339, sliceFrom)
			if err != nil {
				return fmt.Errorf("invalid --from: %w", err)
			}
			opts.TimeFrom = &t
		}
		if sliceTo != "" {
			t, err := time.Parse(time.RFC3339, sliceTo)
			if err != nil {
				return fmt.Errorf("invalid --to: %w", err)
			}
			opts.TimeTo = &t
		}

		result, err := filter.Apply(entries, opts)
		if err != nil {
			return err
		}

		for _, e := range result {
			fmt.Println(e.String())
		}
		return nil
	},
}

func init() {
	sliceCmd.Flags().StringVar(&sliceFrom, "from", "", "start time (RFC3339)")
	sliceCmd.Flags().StringVar(&sliceTo, "to", "", "end time (RFC3339)")
	sliceCmd.Flags().StringVar(&sliceField, "field", "", "field key to match")
	sliceCmd.Flags().StringVar(&sliceFieldVal, "value", "", "field value to match")
	sliceCmd.Flags().StringVar(&sliceTimeField, "time-field", "time", "name of the timestamp field")
	rootCmd.AddCommand(sliceCmd)
}
