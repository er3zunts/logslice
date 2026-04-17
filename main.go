package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "logslice",
	Short: "Filter and slice structured log files by time range or field values",
	Long:  `logslice is a CLI tool to filter structured (JSON) log files by time range or arbitrary field values.`,
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
