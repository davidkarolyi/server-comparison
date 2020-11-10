package job

// Possible states of a benchmark runner
const (
	StateCreated       = "created"
	StateBuildingImage = "building image"
	StateReadyToRun    = "ready to run"
	StateRunning       = "running"
	StateDone          = "done"
)
