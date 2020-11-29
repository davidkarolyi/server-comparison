package reporter

import (
	"fmt"
	"sort"
	"strings"

	"github.com/davidkarolyi/server-comparison/benchmark/chart"
	"github.com/davidkarolyi/server-comparison/benchmark/job"
)

type throughputReporter struct {
	output map[string]float64
}

func (reporter *throughputReporter) Name() string {
	return "throughput"
}

func (reporter *throughputReporter) String() string {
	result := "Throughput:\n"
	for serverName, throughput := range reporter.output {
		result += fmt.Sprintf("  %s: %.2f reqs/sec\n", serverName, throughput)
	}
	return result
}

func (reporter *throughputReporter) Init(measurements map[string]*job.Result) {
	reporter.output = map[string]float64{}
	for fileName, result := range measurements {
		serverName := strings.TrimSuffix(fileName, ".json")
		reporter.output[serverName] = result.WRKResult.ReqsPerSec
	}
}

func (reporter *throughputReporter) ProcessedData() interface{} {
	return reporter.output
}

func (reporter *throughputReporter) Chart() chart.Chart {
	bars := []chart.Bar{}
	for serverName, throughput := range reporter.output {
		bar := chart.Bar{
			Value: throughput,
			Label: serverName,
		}
		bars = append(bars, bar)
	}

	sortBarsDesc(bars)

	return &chart.BarChart{
		Title:     "Throughput [reqs/sec]",
		Precision: 0,
		Bars:      bars,
	}
}

func sortBarsDesc(bars []chart.Bar) {
	sort.SliceStable(bars, func(i, j int) bool {
		return bars[i].Value > bars[j].Value
	})
}
