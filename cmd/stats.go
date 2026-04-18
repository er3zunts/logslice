package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/user/logslice/internal/parser"
	"github.com/user/logslice/internal/stats"
)

var (
	statsFields string
	statsFile   string
)

func init() {
	statsCmd := &cobra.Command{
		Use:   "stats",
		Short: "Print statistics for fields in a log file",
		RunE:  runStats,
	}
	statsCmd.Flags().StringVarP(&statsFile, "file", "f", "", "Input log file (required)")
	statsCmd.Flags().StringVarP(&statsFields, "fields", "F", "level", "Comma-separated fields to summarise")
	_ = statsCmd.MarkFlagRequired("file")
	rootCmd.AddCommand(statsCmd)
}

func runStats(cmd *cobra.Command, args []string) error {
	f, err := os.Open(statsFile)
	if err != nil {
		return fmt.Errorf("open file: %w", err)
	}
	defer f.Close()

	entries, err := parser.ParseJSONLines(f)
	if err != nil {
		return fmt.Errorf("parse: %w", err)
	}

	fields := splitFields(statsFields)
	s := stats.Compute(entries, fields)
	s.Print(os.Stdout)
	return nil
}

func splitFields(s string) []string {
	parts := strings.Split(s, ",")
	out := parts[:0]
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			out = append(out, p)
		}
	}
	return out
}
