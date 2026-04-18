package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/yourorg/logslice/internal/output"
	"github.com/yourorg/logslice/internal/parser"
	"github.com/yourorg/logslice/internal/sample"
)

var (
	sampleHead  int
	sampleTail  int
	sampleEvery int
	sampleRate  float64
	sampleFmt   string
)

func init() {
	sampleCmd := &cobra.Command{
		Use:   "sample [file]",
		Short: "Sample log entries by head, tail, stride, or random rate",
		Args:  cobra.ExactArgs(1),
		RunE:  runSample,
	}
	sampleCmd.Flags().IntVar(&sampleHead, "head", 0, "Keep first N entries")
	sampleCmd.Flags().IntVar(&sampleTail, "tail", 0, "Keep last N entries")
	sampleCmd.Flags().IntVar(&sampleEvery, "every", 0, "Keep every Nth entry")
	sampleCmd.Flags().Float64Var(&sampleRate, "rate", 0, "Random sampling rate (0.0–1.0)")
	sampleCmd.Flags().StringVar(&sampleFmt, "format", "text", "Output format: text, json, pretty")
	rootCmd.AddCommand(sampleCmd)
}

func runSample(cmd *cobra.Command, args []string) error {
	f, err := os.Open(args[0])
	if err != nil {
		return fmt.Errorf("open: %w", err)
	}
	defer f.Close()

	entries, err := parser.ParseJSONLines(f)
	if err != nil {
		return fmt.Errorf("parse: %w", err)
	}

	opts := sample.Options{
		Head:  sampleHead,
		Tail:  sampleTail,
		Every: sampleEvery,
		Rate:  sampleRate,
	}
	result := sample.Apply(entries, opts)

	w, err := output.New(sampleFmt, os.Stdout)
	if err != nil {
		return err
	}
	for _, e := range result {
		if err := w.Write(e); err != nil {
			return err
		}
	}
	return w.Flush()
}
