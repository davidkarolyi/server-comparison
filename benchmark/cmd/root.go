package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{Use: "bmark"}

// Execute runs the constructed cobra rootCommand
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	runBenchmarksCmd.Flags().StringVarP(
		&options.WRKHostURL,
		"remote-url",
		"r",
		"",
		"URL where the remote wrk server is running (required)",
	)
	runBenchmarksCmd.MarkFlagRequired("remote-url")

	runBenchmarksCmd.Flags().StringVarP(
		&options.BenchmarkParams.TargetURL,
		"local-url",
		"l",
		"",
		"URL where the locally hosted server is available (required)",
	)
	runBenchmarksCmd.MarkFlagRequired("local-url")

	runBenchmarksCmd.Flags().IntVarP(
		&options.BenchmarkParams.Connections,
		"connections",
		"c",
		10,
		"Benchmark param: total number of HTTP connections to keep open",
	)

	runBenchmarksCmd.Flags().IntVarP(
		&options.BenchmarkParams.Threads,
		"threads",
		"t",
		2,
		"Benchmark param: total number of threads to use",
	)

	runBenchmarksCmd.Flags().StringVarP(
		&options.BenchmarkParams.Duration,
		"duration",
		"d",
		"10s",
		"Benchmark param: duration of the test, e.g. 2s, 2m, 2h",
	)

	runBenchmarksCmd.Flags().StringVar(
		&options.BenchmarkParams.Timeout,
		"timeout",
		"2s",
		"Benchmark param: record a timeout if a response is not received within this amount of time",
	)

	runBenchmarksCmd.Flags().BoolVarP(
		&options.SkipBuild,
		"skip-build",
		"B",
		false,
		"Prevents benchmarks from building new docker images",
	)

	rootCmd.AddCommand(runWRKServerCmd, runBenchmarksCmd)
}
