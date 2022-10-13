package cmd

import (
	"github.com/rhsyseng/ddosify-tooling/cli-tool/pkg/ddosify"
	"github.com/spf13/cobra"
)

var (
	targetURL    string
	numberOfRuns int
	waitInterval string
)

func addRunFlags(cmd *cobra.Command) {

	flags := cmd.Flags()
	flags.StringVar(&targetURL, "target-url", "", "The target url. e.g: https://google.com")
	flags.IntVar(&numberOfRuns, "runs", 1, "The number of executions.")
	flags.StringVar(&waitInterval, "interval", "1m", "The amount of waiting time between runs.")
	cmd.MarkFlagRequired("target-url")
}

func NewStatusCommand() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "status",
		Short: "Run the status command",
		RunE: func(cmd *cobra.Command, args []string) error {
			lc := ddosify.NewLatencyChecker(targetURL, numberOfRuns, waitInterval)
			return lc.RunCommandStatus()
		},
	}

	addRunFlags(cmd)

	return cmd
}
