package cmd

import (
	"github.com/davidkarolyi/server-comparison/benchmark/wrk"
	"github.com/spf13/cobra"
)

var wrkCmd = &cobra.Command{
	Use:   "wrk",
	Short: "Start wrk server",
	Long:  "Runs a server which is able to run a wrk benchmark against remote URL",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		err := wrk.StartServer()
		return err
	},
}
