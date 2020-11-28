package cmd

import (
	"github.com/davidkarolyi/server-comparison/benchmark/runner"
	"github.com/davidkarolyi/server-comparison/benchmark/wrk/types"
	"github.com/spf13/cobra"
)

var options = &runner.Options{
	WRKHostURL:      "",
	BenchmarkParams: &types.BenchmarkParams{},
}

var runBenchmarksCmd = &cobra.Command{
	Use:   "run",
	Short: "Run server benchmarks",
	Long:  "Runs benchmarks synchronously, using a remote wrk server.",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runner.RunBenchmarks(options)
	},
}
