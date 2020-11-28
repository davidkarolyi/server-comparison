package runner

import (
	"fmt"
	"os"
	"time"

	"github.com/davidkarolyi/server-comparison/benchmark/report"
	"github.com/davidkarolyi/server-comparison/benchmark/wrk/types"
)

// RunBenchmarks will run all server benchmarks, with the help of a remote wrk host.
func RunBenchmarks(wrkHostURL string, benchmarkParams *types.BenchmarkParams) error {
	err := os.Chdir("..")
	if err != nil {
		return err
	}

	jobs, err := newJobList(wrkHostURL, benchmarkParams)
	if err != nil {
		return err
	}
	fmt.Printf("🚦 Servers waiting for benchmark:\n%s\n", jobs)

	err = buildContainers(jobs)
	if err != nil {
		return err
	}

	err = runJobs(jobs)
	if err != nil {
		return err
	}

	err = saveResults(jobs)
	if err != nil {
		return err
	}

	return nil
}

func buildContainers(jobs *jobList) error {
	for _, job := range *jobs {
		fmt.Printf("🏗 Building image for %s...\n", job.ServerName())
		err := job.BuildImages()
		if err != nil {
			return err
		}
	}
	return nil
}

func runJobs(jobs *jobList) error {
	for _, job := range *jobs {
		fmt.Printf("⏱ Benchmarking %s...\n", job.ServerName())
		err := job.Run()
		if err != nil {
			return err
		}
	}
	return nil
}

func saveResults(jobs *jobList) error {
	dirPath := fmt.Sprintf(
		"./reports/%s/_raw/",
		time.Now().Format("2006-01-02"),
	)
	fmt.Printf("📀 Saving results to '%s'", dirPath)
	for _, job := range *jobs {
		result, err := job.Result()
		if err != nil {
			return err
		}
		fileName := fmt.Sprintf("%s.json", job.ServerName())
		report.SaveAsJSON(dirPath+fileName, result)
	}
	return nil
}
