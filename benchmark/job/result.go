package job

import (
	"github.com/davidkarolyi/server-comparison/benchmark/job/docker"
	"github.com/davidkarolyi/server-comparison/benchmark/wrk/types"
)

// Result contains all data collected during a benchmarking job.
type Result struct {
	WRKResult      *types.WRKResult `json:"wrk_result"`
	ContainerStats []*docker.Stats  `json:"container_stats"`
}
