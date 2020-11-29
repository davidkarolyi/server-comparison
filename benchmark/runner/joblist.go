package runner

import (
	"fmt"
	"strings"

	"github.com/davidkarolyi/server-comparison/benchmark/job"
	"github.com/davidkarolyi/server-comparison/benchmark/utils"
	"github.com/davidkarolyi/server-comparison/benchmark/wrk"
)

// jobList is a list of benchmark jobs
type jobList []*job.BenchmarkJob

// newJobList will create a job list from server directories found in the project root.
func newJobList(options *Options) (*jobList, error) {
	_, dirNames, err := utils.ListDirContent("./")
	if err != nil {
		return nil, err
	}

	wrkClient := wrk.NewClient(options.WRKHostURL, options.BenchmarkParams)
	err = wrkClient.CheckConnection()
	if err != nil {
		return nil, err
	}

	jobs := &jobList{}
	for _, dirName := range dirNames {
		if isServerDir(dirName) {
			runner, err := job.New(dirName, wrkClient)
			if err != nil {
				return nil, err
			}
			jobs.push(runner)
		}
	}
	return jobs, nil
}

func (jobs *jobList) String() string {
	stringForm := ""
	for _, job := range *jobs {
		stringForm += fmt.Sprintln("-", job.ServerName())
	}
	return stringForm
}

func (jobs *jobList) push(runner *job.BenchmarkJob) {
	*jobs = append(*jobs, runner)
}

func isServerDir(dirName string) bool {
	return !isHiddenDir(dirName) && hasUnderscoreInName(dirName)
}

func isHiddenDir(dirName string) bool {
	return strings.HasPrefix(dirName, ".")
}

func hasUnderscoreInName(dirName string) bool {
	return strings.Contains(dirName, "_")
}
