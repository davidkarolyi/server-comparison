package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{Use: "bmark"}

// Execute runs the constructed cobra rootCommand
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	initFlagsRunCmd()
	initFlagsReportCmd()

	rootCmd.AddCommand(
		wrkCmd,
		runCmd,
		reportCmd,
	)
}
