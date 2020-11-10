package wrk

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"

	"github.com/davidkarolyi/server-comparison/benchmark/wrk/convert"
	"github.com/davidkarolyi/server-comparison/benchmark/wrk/types"
)

// wrkRunner runs wrk via docker: https://hub.docker.com/r/williamyeh/wrk/
type wrkRunner struct {
	initialized bool
}

// newWRKRunner creates a new wrkRunner instance
func newWRKRunner() *wrkRunner {
	return &wrkRunner{
		initialized: false,
	}
}

// Init pulls the wrk container image
func (runner *wrkRunner) Init() error {
	runner.initialized = true

	pullCommand := "docker pull williamyeh/wrk"
	cmd := exec.Command("bash", "-c", pullCommand)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout

	return cmd.Run()
}

// RunBenchmark runs a wrk benchmark from a docker container
func (runner *wrkRunner) RunBenchmark(ctx context.Context, params *types.BenchmarkParams) (*types.WRKResult, error) {
	if !runner.initialized {
		return nil, errors.New("Benchmark runner hasn't initialised yet, call wrkRunner.Init() to do so")
	}

	outputBuffer := new(bytes.Buffer)

	cmd := exec.CommandContext(ctx, "bash", "-c", buildBenchmarkCommand(params))
	cmd.Stderr = os.Stderr
	cmd.Stdout = outputBuffer

	err := cmd.Run()
	if err != nil {
		return nil, err
	}

	return convert.TerminalOutputToResult(outputBuffer)
}

func buildBenchmarkCommand(params *types.BenchmarkParams) string {
	flags := fmt.Sprintf(
		"-t%d -c%d -d%s --timeout %s",
		params.Threads,
		params.Connections,
		params.Duration,
		params.Timeout,
	)
	return fmt.Sprintf(
		"docker run --rm  --network=host -v `pwd`:/data williamyeh/wrk %s %s",
		params.TargetURL,
		flags,
	)
}
