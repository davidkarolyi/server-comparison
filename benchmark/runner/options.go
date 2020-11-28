package runner

import "github.com/davidkarolyi/server-comparison/benchmark/wrk/types"

// Options contains configurations for running the benchmark suite.
type Options struct {
	WRKHostURL      string
	BenchmarkParams *types.BenchmarkParams
}
