package runner

import "github.com/davidkarolyi/server-comparison/benchmark/wrk/types"

// Options contains configurations for running the benchmark suite.
type Options struct {
	WRKHostURL      string                 `json:"wrk_host_url"`
	BenchmarkParams *types.BenchmarkParams `json:"benchmark_params"`
	SkipBuild       bool                   `json:"skip_build"`
}
