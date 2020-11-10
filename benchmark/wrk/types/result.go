package types

// WRKResult houses wrk related benchmark results.
type WRKResult struct {
	RawOutput             string       `json:"raw_output"`
	ReqsPerSec            float64      `json:"reqs_per_sec"`
	TransferedBytesPerSec int64        `json:"transfered_bytes_per_sec"`
	ThreadStats           *ThreadStats `json:"thread_stats"`
}

// ThreadStats holds numbers obtained per thread.
// Learn more: https://github.com/wg/wrk/issues/259.
type ThreadStats struct {
	Latency   *Latency   `json:"latency"`
	ReqPerSec *ReqPerSec `json:"req_per_sec"`
}

// Latency holds request latency related statistics.
type Latency struct {
	AvgMS           float64 `json:"avg_ms"`
	StdevMS         float64 `json:"stdev_ms"`
	MaxMS           float64 `json:"max_ms"`
	ReqsInStdevPerc float64 `json:"reqs_in_stdev_perc"`
}

// ReqPerSec holds requests per seconds statistics.
type ReqPerSec struct {
	Avg             float64 `json:"avg"`
	Stdev           float64 `json:"stdev"`
	Max             float64 `json:"max"`
	ReqsInStdevPerc float64 `json:"reqs_in_stdev_perc"`
}
