package runner

import (
	"fmt"
	"time"

	"github.com/davidkarolyi/server-comparison/benchmark/utils"
)

// RunBenchmarks will run all server benchmarks, with the help of a remote wrk host.
func RunBenchmarks(options *Options) error {

	err := utils.ChangeToProjectRoot()
	if err != nil {
		return err
	}

	jobs, err := newJobList(options)
	if err != nil {
		return err
	}
	fmt.Printf("🚦 Servers waiting for benchmark:\n%s\n", jobs)

	if options.SkipBuild {
		skipBuilds(jobs)
	} else {
		err = buildContainers(jobs)
		if err != nil {
			return err
		}
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

func skipBuilds(jobs *jobList) {
	for _, job := range *jobs {
		fmt.Printf("⏩ Skipping build process of %s...\n", job.ServerName())
		job.SkipBuild()
	}
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
		time.Now().Format(time.RFC3339),
	)
	fmt.Printf("📀 Saving results to '%s'", dirPath)
	for _, job := range *jobs {
		result, err := job.Result()
		if err != nil {
			return err
		}
		fileName := fmt.Sprintf("%s.json", job.ServerName())
		utils.SaveAsJSON(dirPath+fileName, result)
	}
	return nil
}
