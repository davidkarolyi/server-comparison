package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{Use: "benchmark"}

// Execute runs the constructed cobra rootCommand
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	runBenchmarksCmd.Flags().StringVarP(
		&localHostURL,
		"local",
		"l",
		"",
		"URL where the locally hosted server is available (required)",
	)
	runBenchmarksCmd.MarkFlagRequired("local")

	runBenchmarksCmd.Flags().StringVarP(
		&wrkHostURL,
		"remote",
		"r",
		"",
		"URL where the remote wrk server is running (required)",
	)
	runBenchmarksCmd.MarkFlagRequired("remote")
	rootCmd.AddCommand(runWRKServerCmd, runBenchmarksCmd)
}
