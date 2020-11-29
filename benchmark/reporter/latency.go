package reporter

import (
	"fmt"
	"strings"

	"github.com/davidkarolyi/server-comparison/benchmark/chart"
	"github.com/davidkarolyi/server-comparison/benchmark/job"
)

type latencyReporter struct {
	output map[string]float64
}

func (reporter *latencyReporter) Name() string {
	return "average_latency"
}

func (reporter *latencyReporter) String() string {
	result := "Average Latency:\n"
	for serverName, latency := range reporter.output {
		result += fmt.Sprintf("  %s: %.2f ms\n", serverName, latency)
	}
	return result
}

func (reporter *latencyReporter) Init(measurements map[string]*job.Result) {
	reporter.output = map[string]float64{}
	for fileName, result := range measurements {
		serverName := strings.TrimSuffix(fileName, ".json")
		reporter.output[serverName] = result.WRKResult.ThreadStats.Latency.AvgMS
	}
}

func (reporter *latencyReporter) ProcessedData() interface{} {
	return reporter.output
}

func (reporter *latencyReporter) Chart() chart.Chart {
	bars := []chart.Bar{}
	for serverName, latency := range reporter.output {
		bar := chart.Bar{
			Value: latency,
			Label: serverName,
		}
		bars = append(bars, bar)
	}

	sortBarsAsc(bars)

	return &chart.BarChart{
		Title:     "Average Latency [ms]",
		Precision: 2,
		Bars:      bars,
	}
}
