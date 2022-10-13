package cmd

import (
	"errors"

	"github.com/rhsyseng/ddosify-tooling/cli-tool/pkg/ddosify"
	"github.com/spf13/cobra"
)

var (
	targetURL    string
	numberOfRuns int
	waitInterval string
	locations    []string
)

func addExecFlags(cmd *cobra.Command) {

	flags := cmd.Flags()
	flags.StringVarP(&targetURL, "target-url", "t", "", "The target url. e.g: https://google.com")
	flags.IntVarP(&numberOfRuns, "runs", "r", 1, "The number of executions.")
	flags.StringVarP(&waitInterval, "interval", "i", "1m", "The amount of waiting time between runs.")
	flags.StringArrayVarP(&locations, "locations", "l", []string{"EU.ES.*"}, "The array of locations to be requested. e.g: NA.US.*,NA.EU.*")
	cmd.MarkFlagRequired("target-url")
}

func NewExecCommand() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "run",
		Short: "Exec the run command",
		RunE: func(cmd *cobra.Command, args []string) error {
			// Validate URL is valid
			validURL := ddosify.ValidateURL(targetURL)
			if !validURL {
				return errors.New(" not valid url")
			}
			lc := ddosify.NewLatencyChecker(targetURL, numberOfRuns, waitInterval, locations)
			return lc.RunCommandExec()
		},
	}

	addExecFlags(cmd)

	return cmd
}
