package convert_test

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/davidkarolyi/server-comparison/benchmark/wrk/convert"
	"github.com/davidkarolyi/server-comparison/benchmark/wrk/types"
	"github.com/google/go-cmp/cmp"
)

var outputBuffer = bytes.NewBufferString(
	`Running 10s test @ http://localhost:3000
  2 threads and 10 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
	Latency    11.10ms   22.90ms 262.28ms   97.74%
	Req/Sec   624.00    179.05     0.98k    68.75%
  12009 requests in 10.05s, 3.56MB read
  Socket errors: connect 0, read 1474, write 0, timeout 0
	Requests/sec:   1194.72
	Transfer/sec:    362.85KB`,
)

func TestConvertOutputToResult(t *testing.T) {
	t.Run("returns error on invalid input", func(t *testing.T) {
		_, err := convert.TerminalOutputToResult(bytes.NewBufferString("foo bar"))
		if err == nil {
			t.Fatal("Error is nil, but expected an actual value")
		}
	})

	t.Run("returns no error on valid input", func(t *testing.T) {
		result, err := convert.TerminalOutputToResult(outputBuffer)
		if result == nil || err != nil {
			t.Fatal(err)
		}
	})

	t.Run("returns the raw wrk output", func(t *testing.T) {
		result, err := convert.TerminalOutputToResult(outputBuffer)
		if err != nil {
			t.Fatal(err)
		}
		expected := outputBuffer.String()
		actual := result.RawOutput
		if actual != expected {
			t.Fatalf("Expected: '%s', Got: '%s'", expected, actual)
		}
	})

	t.Run("returns correct ReqsPerSec value", func(t *testing.T) {
		result, err := convert.TerminalOutputToResult(outputBuffer)
		if err != nil {
			t.Fatal(err)
		}
		expected := 1194.72
		actual := result.ReqsPerSec
		if actual != expected {
			t.Fatalf("Expected: %f, Got: %f", expected, actual)
		}
	})

	t.Run("returns correct TransferedBytesPerSec value", func(t *testing.T) {
		result, err := convert.TerminalOutputToResult(outputBuffer)
		if err != nil {
			t.Fatal(err)
		}

		var expected int64 = 371558 // 362.85 * 1024
		actual := result.TransferedBytesPerSec
		if actual != expected {
			t.Fatalf("Expected: %d, Got: %d", expected, actual)
		}
	})

	t.Run("returns correct Latency.AvgMS value", func(t *testing.T) {
		result, err := convert.TerminalOutputToResult(outputBuffer)
		if err != nil {
			t.Fatal(err)
		}
		expected := 11.10
		actual := result.ThreadStats.Latency.AvgMS
		if actual != expected {
			t.Fatalf("Expected: %fms, Got: %fms", expected, actual)
		}
	})

	t.Run("returns correct Latency.MaxMS value", func(t *testing.T) {
		result, err := convert.TerminalOutputToResult(outputBuffer)
		if err != nil {
			t.Fatal(err)
		}
		expected := 262.28
		actual := result.ThreadStats.Latency.MaxMS
		if actual != expected {
			t.Fatalf("Expected: %fms, Got: %fms", expected, actual)
		}
	})

	t.Run("returns correct Latency.StdevMS value", func(t *testing.T) {
		result, err := convert.TerminalOutputToResult(outputBuffer)
		if err != nil {
			t.Fatal(err)
		}
		expected := 22.90
		actual := result.ThreadStats.Latency.StdevMS
		if actual != expected {
			t.Fatalf("Expected: %fms, Got: %fms", expected, actual)
		}
	})

	t.Run("returns correct Latency.ReqsInStdevPerc value", func(t *testing.T) {
		result, err := convert.TerminalOutputToResult(outputBuffer)
		if err != nil {
			t.Fatal(err)
		}
		expected := 97.74
		actual := result.ThreadStats.Latency.ReqsInStdevPerc
		if actual != expected {
			t.Fatalf("Expected: %f%%, Got: %f%%", expected, actual)
		}
	})

	t.Run("returns correct ReqPerSec.Avg value", func(t *testing.T) {
		result, err := convert.TerminalOutputToResult(outputBuffer)
		if err != nil {
			t.Fatal(err)
		}
		expected := 624.00
		actual := result.ThreadStats.ReqPerSec.Avg
		if actual != expected {
			t.Fatalf("Expected: %f, Got: %f", expected, actual)
		}
	})

	t.Run("returns correct ReqPerSec.Max value", func(t *testing.T) {
		result, err := convert.TerminalOutputToResult(outputBuffer)
		if err != nil {
			t.Fatal(err)
		}
		expected := 980.0
		actual := result.ThreadStats.ReqPerSec.Max
		if actual != expected {
			t.Fatalf("Expected: %f, Got: %f", expected, actual)
		}
	})

	t.Run("returns correct ReqPerSec.Stdev value", func(t *testing.T) {
		result, err := convert.TerminalOutputToResult(outputBuffer)
		if err != nil {
			t.Fatal(err)
		}
		expected := 179.05
		actual := result.ThreadStats.ReqPerSec.Stdev
		if actual != expected {
			t.Fatalf("Expected: %f, Got: %f", expected, actual)
		}
	})

	t.Run("returns correct ReqPerSec.ReqsInStdevPerc value", func(t *testing.T) {
		result, err := convert.TerminalOutputToResult(outputBuffer)
		if err != nil {
			t.Fatal(err)
		}
		expected := 68.75
		actual := result.ThreadStats.ReqPerSec.ReqsInStdevPerc
		if actual != expected {
			t.Fatalf("Expected: %f, Got: %f", expected, actual)
		}
	})

	t.Run("functions properly if not all line was printed", func(t *testing.T) {
		outputBufferWithMissingLines := bytes.NewBufferString(
			`Running 10s test @ http://localhost:3000
				Thread Stats   Avg      Stdev     Max   +/- Stdev
					Latency     9.41ms    4.25ms  73.92ms   88.68%
					Req/Sec   536.23    140.78   818.00     65.00%
			Requests/sec:   1067.02
			Transfer/sec:    324.06KB`,
		)

		expected := &types.WRKResult{
			RawOutput:             outputBufferWithMissingLines.String(),
			ReqsPerSec:            1067.02,
			TransferedBytesPerSec: 331837, // 324.06 * 1024
			ThreadStats: &types.ThreadStats{
				Latency: &types.Latency{
					AvgMS:           9.41,
					StdevMS:         4.25,
					MaxMS:           73.92,
					ReqsInStdevPerc: 88.68,
				},
				ReqPerSec: &types.ReqPerSec{
					Avg:             536.23,
					Stdev:           140.78,
					Max:             818,
					ReqsInStdevPerc: 65,
				},
			},
		}

		actual, err := convert.TerminalOutputToResult(outputBufferWithMissingLines)
		if err != nil {
			t.Fatal(err)
		}

		if !cmp.Equal(actual, expected) {
			t.Fatalf("Expected:\n%+v\nGot:\n%+v\n", prettyPrint(expected), prettyPrint(actual))
		}
	})
}

func prettyPrint(result *types.WRKResult) string {
	resultJSON, _ := json.MarshalIndent(result, "", "  ")
	return string(resultJSON)
}
