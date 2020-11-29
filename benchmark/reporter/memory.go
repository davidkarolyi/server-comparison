package reporter

import (
	"fmt"
	"math"
	"sort"
	"strings"

	"github.com/davidkarolyi/server-comparison/benchmark/chart"
	"github.com/davidkarolyi/server-comparison/benchmark/job"
)

type memoryUsageReporter struct {
	output map[string]float64
}

func (reporter *memoryUsageReporter) Name() string {
	return "peak_memory_usage_in_mb"
}

func (reporter *memoryUsageReporter) String() string {
	result := "Memory Usage:\n"
	for serverName, memoryUsage := range reporter.output {
		result += fmt.Sprintf("  %s: %.2f MB\n", serverName, memoryUsage)
	}
	return result
}

func (reporter *memoryUsageReporter) Init(measurements map[string]*job.Result) {
	reporter.output = map[string]float64{}
	for fileName, result := range measurements {
		serverName := strings.TrimSuffix(fileName, ".json")
		reporter.output[serverName] = bytesToMegaBytes(getPeakMemoryUsage(result))
	}
}

func (reporter *memoryUsageReporter) ProcessedData() interface{} {
	return reporter.output
}

func (reporter *memoryUsageReporter) Chart() chart.Chart {
	bars := []chart.Bar{}
	for serverName, memoryUsage := range reporter.output {
		bar := chart.Bar{
			Value: memoryUsage,
			Label: serverName,
		}
		bars = append(bars, bar)
	}

	sortBarsAsc(bars)

	return &chart.BarChart{
		Title:     "Peak Memory Usage [MB]",
		Precision: 2,
		Bars:      bars,
	}
}

func getPeakMemoryUsage(result *job.Result) (peak uint64) {
	for _, statsSnapshot := range result.ContainerStats {
		usage := statsSnapshot.MemoryStats.MaxUsage
		if usage > peak {
			peak = usage
		}
	}
	return
}

func bytesToMegaBytes(bytes uint64) float64 {
	return float64(bytes) / math.Pow(1024, 2)
}

func sortBarsAsc(bars []chart.Bar) {
	sort.SliceStable(bars, func(i, j int) bool {
		return bars[i].Value < bars[j].Value
	})
}
