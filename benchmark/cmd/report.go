package cmd

import (
	"github.com/davidkarolyi/server-comparison/benchmark/reporter"
	"github.com/davidkarolyi/server-comparison/benchmark/utils"
	"github.com/spf13/cobra"
)

var preview bool
var source string

var reportCmd = &cobra.Command{
	Use:   "report",
	Short: "Generates report from the latest benchmark",
	Long:  "Processes the data in the 'measurements' directory and generates SVG charts from them.",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		err := utils.ChangeToProjectRoot()
		if err != nil {
			return err
		}

		if source == "" {
			source, err = reporter.LatestReportName()
			if err != nil {
				return err
			}
		}

		report, err := reporter.NewReport(source)
		if err != nil {
			return err
		}

		if preview {
			report.Preview()
			return nil
		}
		return report.Generate()
	},
}

func initFlagsReportCmd() {
	reportCmd.Flags().BoolVarP(
		&preview,
		"preview",
		"p",
		false,
		"Show a preview of the report, without saving it to disk",
	)

	reportCmd.Flags().StringVarP(
		&source,
		"source",
		"s",
		"",
		"Show a preview of the report, without saving it to disk",
	)
}
