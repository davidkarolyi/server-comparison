package job

import (
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/davidkarolyi/server-comparison/benchmark/job/docker"
	"github.com/davidkarolyi/server-comparison/benchmark/wrk"
)

// BenchmarkJob can run a benchmark for a specific server.
type BenchmarkJob struct {
	serverName      string
	state           string
	docker          docker.IDocker
	wrkClient       *wrk.Client
	statsRecorder   *statsRecorder
	cleanedUpSignal chan struct{}
	result          *Result
}

// New creates a new benchmark job.
func New(serverName string, wrkClient *wrk.Client) (*BenchmarkJob, error) {
	dockerService, err := docker.New(serverName)
	if err != nil {
		return nil, err
	}

	job := &BenchmarkJob{
		serverName:      serverName,
		state:           StateCreated,
		docker:          dockerService,
		wrkClient:       wrkClient,
		statsRecorder:   newStatsRecorder(),
		cleanedUpSignal: make(chan struct{}),
	}

	job.setupGracefullShutdown()

	return job, nil
}

// State returns the state of the actual job.
func (job *BenchmarkJob) State() string {
	return job.state
}

// ServerName returns the name of the server, which is associated with the actual job.
func (job *BenchmarkJob) ServerName() string {
	return job.serverName
}

// Result will return the benchmarking job's result if it's available
func (job *BenchmarkJob) Result() (*Result, error) {
	if job.result == nil {
		return nil, errors.New("Result is not available")
	}
	return job.result, nil
}

// BuildImages builds necesarry images for the benchmark.
func (job *BenchmarkJob) BuildImages() error {
	if job.state != StateCreated {
		return fmt.Errorf("Cannot build image, just after %s state", StateCreated)
	}
	job.state = StateBuildingImage

	fmt.Printf("🏗 Building image for %s...\n", job.ServerName())
	err := job.docker.BuildServerImage()
	if err != nil {
		return fmt.Errorf("Cannot build image for %s: %s", job.serverName, err)
	}

	job.state = StateReadyToRun
	return nil
}

// SkipBuild will skip the build step for this job
func (job *BenchmarkJob) SkipBuild() {
	if job.state == StateCreated {
		fmt.Printf("⏩ Skipping build process of %s...\n", job.ServerName())
		job.state = StateReadyToRun
	}
}

// Run will run the benchmark for the actual server.
func (job *BenchmarkJob) Run() error {
	if job.state != StateReadyToRun {
		return fmt.Errorf("Cannot run image, just after %s state", StateReadyToRun)
	}

	job.state = StateRunning
	defer job.terminate()

	fmt.Printf("⏱ Starting conatiner for %s...\n", job.ServerName())
	statsChannel, err := job.docker.StartServerContainer()
	if err != nil {
		return err
	}

	err = job.statsRecorder.AddStatsChannel(statsChannel)
	if err != nil {
		return err
	}

	fmt.Println("🧘‍♀️ Waiting 5 seconds to make sure the server is in an idle state...")
	time.Sleep(5 * time.Second)

	err = job.statsRecorder.Record()
	if err != nil {
		return err
	}

	fmt.Println("🚀 Running wrk benchmark...")
	wrkResult, err := job.wrkClient.RunBenchmark()
	if err != nil {
		return err
	}

	containerStats, err := job.statsRecorder.Done()
	if err != nil {
		return err
	}

	job.state = StateDone
	job.result = &Result{
		WRKResult:      wrkResult,
		ContainerStats: containerStats,
	}
	return nil
}

func (job *BenchmarkJob) setupGracefullShutdown() {
	signals := make(chan os.Signal, 1)

	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		select {
		case <-signals:
			fmt.Println()
			job.terminate()
			os.Exit(1)
		case <-job.cleanedUpSignal:
			return
		}
	}()
}

func (job *BenchmarkJob) terminate() {
	fmt.Println("✨ Cleaning up running containers...")
	job.docker.Clean()
	job.cleanedUpSignal <- struct{}{}
}
