package convert

import (
	"bytes"

	"github.com/davidkarolyi/server-comparison/benchmark/wrk/types"
)

// TerminalOutputToResult converts wrk's stdout into a Result object.
func TerminalOutputToResult(outputBuffer *bytes.Buffer) (*types.WRKResult, error) {
	matrix := newValueMatrix(outputBuffer.String())
	extractor := newFieldExtractor(matrix)

	result := newResult()
	var err error = nil

	result.RawOutput = outputBuffer.String()

	result.ReqsPerSec, err = extractor.ReqsPerSec()
	if err != nil {
		return nil, err
	}

	result.TransferedBytesPerSec, err = extractor.TransferedBytesPerSec()
	if err != nil {
		return nil, err
	}

	result.ThreadStats.Latency.AvgMS, err = extractor.LatencyAvgMS()
	if err != nil {
		return nil, err
	}

	result.ThreadStats.Latency.MaxMS, err = extractor.LatencyMaxMS()
	if err != nil {
		return nil, err
	}

	result.ThreadStats.Latency.StdevMS, err = extractor.LatencyStdevMS()
	if err != nil {
		return nil, err
	}

	result.ThreadStats.Latency.ReqsInStdevPerc, err = extractor.LatencyReqsInStdevPerc()
	if err != nil {
		return nil, err
	}

	result.ThreadStats.ReqPerSec.Avg, err = extractor.ReqPerSecAvg()
	if err != nil {
		return nil, err
	}

	result.ThreadStats.ReqPerSec.Max, err = extractor.ReqPerSecMax()
	if err != nil {
		return nil, err
	}

	result.ThreadStats.ReqPerSec.Stdev, err = extractor.ReqPerSecStdev()
	if err != nil {
		return nil, err
	}

	result.ThreadStats.ReqPerSec.ReqsInStdevPerc, err = extractor.ReqPerSecReqsInStdevPerc()
	if err != nil {
		return nil, err
	}

	return result, nil
}

func newResult() *types.WRKResult {
	return &types.WRKResult{
		ThreadStats: &types.ThreadStats{
			Latency:   &types.Latency{},
			ReqPerSec: &types.ReqPerSec{},
		},
	}
}
