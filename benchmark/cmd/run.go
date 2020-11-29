package cmd

import (
	"github.com/davidkarolyi/server-comparison/benchmark/runner"
	"github.com/davidkarolyi/server-comparison/benchmark/utils"
	"github.com/davidkarolyi/server-comparison/benchmark/wrk/types"
	"github.com/spf13/cobra"
)

var options = &runner.Options{
	WRKHostURL:      "",
	BenchmarkParams: &types.BenchmarkParams{},
	SkipBuild:       false,
}

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run server benchmarks",
	Long:  "Runs benchmarks synchronously, using a remote wrk server.",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		err := utils.ChangeToProjectRoot()
		if err != nil {
			return err
		}

		return runner.RunBenchmarks(options)
	},
}

func initFlagsRunCmd() {
	runCmd.Flags().StringVarP(
		&options.WRKHostURL,
		"remote-url",
		"r",
		"",
		"URL where the remote wrk server is running (required)",
	)
	runCmd.MarkFlagRequired("remote-url")

	runCmd.Flags().StringVarP(
		&options.BenchmarkParams.TargetURL,
		"local-url",
		"l",
		"",
		"URL where the locally hosted server is available (required)",
	)
	runCmd.MarkFlagRequired("local-url")

	runCmd.Flags().IntVarP(
		&options.BenchmarkParams.Connections,
		"connections",
		"c",
		100,
		"Benchmark param: total number of HTTP connections to keep open",
	)

	runCmd.Flags().IntVarP(
		&options.BenchmarkParams.Threads,
		"threads",
		"t",
		2,
		"Benchmark param: total number of threads to use",
	)

	runCmd.Flags().StringVarP(
		&options.BenchmarkParams.Duration,
		"duration",
		"d",
		"30s",
		"Benchmark param: duration of the test, e.g. 2s, 2m, 2h",
	)

	runCmd.Flags().StringVar(
		&options.BenchmarkParams.Timeout,
		"timeout",
		"2s",
		"Benchmark param: record a timeout if a response is not received within this amount of time",
	)

	runCmd.Flags().BoolVarP(
		&options.SkipBuild,
		"skip-build",
		"B",
		false,
		"Prevents benchmarks from building new docker images",
	)
}
