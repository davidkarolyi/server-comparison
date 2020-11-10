package runner

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/davidkarolyi/server-comparison/benchmark/job"
	"github.com/davidkarolyi/server-comparison/benchmark/wrk"
)

// jobList is a list of benchmark jobs
type jobList []*job.BenchmarkJob

// newJobList will create a job list from server directories found in the project root.
func newJobList(localHostURL string, wrkHostURL string) (*jobList, error) {
	files, err := ioutil.ReadDir("./")
	if err != nil {
		return nil, err
	}

	wrkClient := wrk.NewClient(wrkHostURL, wrk.DefaultParamsWithTargetURL(localHostURL))
	err = wrkClient.CheckConnection()
	if err != nil {
		return nil, err
	}

	jobs := &jobList{}
	for _, file := range files {
		if isServerDir(file) {
			runner, err := job.New(file.Name(), wrkClient)
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

func isServerDir(file os.FileInfo) bool {
	return file.IsDir() && !isHiddenDir(file) && hasUnderscoreInName(file)
}

func isHiddenDir(file os.FileInfo) bool {
	return strings.HasPrefix(file.Name(), ".")
}

func hasUnderscoreInName(file os.FileInfo) bool {
	return strings.Contains(file.Name(), "_")
}
