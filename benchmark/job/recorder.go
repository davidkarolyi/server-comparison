package job

import (
	"errors"

	"github.com/davidkarolyi/server-comparison/benchmark/job/docker"
)

// statsRecorder records container statistics of a running container.
type statsRecorder struct {
	started        bool
	final          bool
	containerStats []*docker.Stats
	statsChannel   chan *docker.Stats
}

// New creates initializes a new recording.
func newStatsRecorder() *statsRecorder {
	return &statsRecorder{
		started:        false,
		final:          false,
		statsChannel:   nil,
		containerStats: []*docker.Stats{},
	}
}

// AddStatsChannel will initialize monitoring of the given stats channel.
func (rec *statsRecorder) AddStatsChannel(statsChannel chan *docker.Stats) error {
	if rec.statsChannel != nil {
		return errors.New("Recorder has already got a source")
	}
	rec.statsChannel = statsChannel
	go rec.monitorStats()
	return nil
}

// Record starts the recording.
func (rec *statsRecorder) Record() error {
	if !rec.final {
		rec.started = true
		return nil
	}
	return errors.New("Cannot restart a stopped recorder")
}

// Done stops the recording.
func (rec *statsRecorder) Done() ([]*docker.Stats, error) {
	if rec.final {
		return nil, errors.New("Recorder has already stopped")
	}
	if !rec.started {
		return nil, errors.New("Recorder was not started yet")
	}
	rec.final = true

	return rec.containerStats, nil
}

func (rec *statsRecorder) monitorStats() {
	for {
		snapshot, ok := <-rec.statsChannel
		if !ok {
			break
		}
		if rec.started && !rec.final {
			rec.containerStats = append(rec.containerStats, snapshot)
		}
	}
}
