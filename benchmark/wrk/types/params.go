package types

// BenchmarkParams contains information about the requested wrk benchmark.
type BenchmarkParams struct {
	TargetURL   string `json:"target_url"`
	Duration    string `json:"duration"`
	Threads     int    `json:"threads"`
	Connections int    `json:"connections"`
	Timeout     string `json:"timeout"`
}
