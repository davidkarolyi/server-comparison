package docker

import (
	"encoding/json"

	"github.com/docker/docker/api/types"
)

// Stats represents a snapshot of resource utilization of a running continer.
type Stats types.Stats

// statsWriter implements the io.Writer interface,
// and writes stats objects into the given channel.
type statsWriter struct {
	statsChannel chan *Stats
}

func newStatsWriter(statsChannel chan *Stats) *statsWriter {
	return &statsWriter{
		statsChannel: statsChannel,
	}
}

func (writer *statsWriter) Write(p []byte) (n int, err error) {
	stats := &Stats{}
	err = json.Unmarshal(p, stats)
	if err != nil {
		return 0, err
	}

	writer.statsChannel <- stats

	return len(p), nil
}
