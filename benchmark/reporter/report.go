package reporter

import (
	"errors"
	"fmt"
	"os"

	"github.com/davidkarolyi/server-comparison/benchmark/chart"
	"github.com/davidkarolyi/server-comparison/benchmark/job"
	"github.com/davidkarolyi/server-comparison/benchmark/utils"
)

// reporter can produce a processed data output and a chart from measurements.
type reporter interface {
	Name() string
	String() string
	Init(measurements map[string]*job.Result)
	ProcessedData() interface{}
	Chart() chart.Chart
}

var reporters = []reporter{
	new(throughputReporter),
	new(memoryUsageReporter),
	new(latencyReporter),
}

// LatestReportName will return the name of the latest report
// can be found in the reports directory.
func LatestReportName() (string, error) {
	_, dirNames, err := utils.ListDirContent("./reports")
	if err != nil {
		return "", err
	}

	if len(dirNames) == 0 {
		return "", errors.New("There are no reports yet")
	}

	lastIndex := len(dirNames) - 1
	return dirNames[lastIndex], nil
}

// Report represents all data related to a report
type Report struct {
	pathToReport string
	Measurements map[string]*job.Result
	Reporters    []reporter
}

// NewReport will load measurement data from the given report folder,
// and creates a new report object from it. e.g.: source="2020-11-29T13:55:27Z".
func NewReport(source string) (*Report, error) {
	report := &Report{
		Measurements: map[string]*job.Result{},
		Reporters:    reporters,
		pathToReport: fmt.Sprintf("./reports/%s", source),
	}

	err := report.populateMeasurements()
	if err != nil {
		return nil, err
	}

	return report, nil
}

// Generate generates the actual report.
func (report *Report) Generate() error {
	for _, reporter := range report.Reporters {
		fmt.Printf("📎 Generating %s report\n", reporter.Name())
		reporter.Init(report.Measurements)

		err := report.runReporter(reporter)
		if err != nil {
			return err
		}
	}

	fmt.Printf("✅ Report generated under %s\n", report.pathToReport)
	return nil
}

// Preview will print a preview of the processed data.
func (report *Report) Preview() {
	for _, reporter := range report.Reporters {
		reporter.Init(report.Measurements)
		fmt.Println(reporter)
	}
}

func (report *Report) populateMeasurements() error {
	pathToMeasurements := fmt.Sprintf("%s/measurements", report.pathToReport)
	fmt.Printf("📖 Reading measurement data from %s\n", pathToMeasurements)

	fileNames, _, err := utils.ListDirContent(pathToMeasurements)
	if err != nil {
		return err
	}

	for _, fileName := range fileNames {
		report.Measurements[fileName] = &job.Result{}
		err := utils.ReadJSON(
			fmt.Sprintf("%s/%s", pathToMeasurements, fileName),
			report.Measurements[fileName],
		)
		if err != nil {
			return err
		}
	}

	return nil
}

func (report *Report) runReporter(reporter reporter) error {
	dataPath := fmt.Sprintf("%s/%s.json", report.pathToReport, reporter.Name())
	chartPath := fmt.Sprintf("%s/%s.svg", report.pathToReport, reporter.Name())

	fmt.Printf("📝 Precessing %s data\n", reporter.Name())
	data := reporter.ProcessedData()
	err := utils.SaveAsJSON(dataPath, data)
	if err != nil {
		return err
	}

	fmt.Printf("📊 Generating %s chart\n", reporter.Name())
	chartFile, err := os.Create(chartPath)
	if err != nil {
		return err
	}
	defer chartFile.Close()

	reporter.Chart().Render(chartFile)

	return nil
}
