package cmd

import (
	"errors"
	"log"

	"github.com/rhsyseng/ddosify-tooling/cli-tool/pkg/ddosify"
	"github.com/spf13/cobra"
)

var (
	targetURL             string
	numberOfRuns          int
	waitInterval          string
	locations             []string
	outputLocationsNumber int
)

func addExecFlags(cmd *cobra.Command) {

	flags := cmd.Flags()
	flags.StringVarP(&targetURL, "target-url", "t", "", "The target url. e.g: https://google.com")
	flags.IntVarP(&numberOfRuns, "runs", "r", 1, "The number of executions.")
	flags.StringVarP(&waitInterval, "interval", "i", "1m", "The amount of waiting time between runs.")
	flags.StringArrayVarP(&locations, "locations", "l", []string{"EU.ES.*"}, "The array of locations to be requested. e.g: NA.US.*,NA.EU.*")
	flags.IntVar(&outputLocationsNumber, "output-locations", 1, "The number of best locations to output.")
	cmd.MarkFlagRequired("target-url")
}

func NewExecCommand() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "run",
		Short: "Exec the run command",
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: Create validator function for all flags
			// TODO: Validate interval is valid
			validInterval := ddosify.ValidateIntervalTime(waitInterval)
			if !validInterval {
				return errors.New(" not valid interval")
			}
			// Validate URL is valid
			validURL := ddosify.ValidateURL(targetURL)
			if !validURL {
				return errors.New(" not valid url")
			}
			// Get waitIntervalInSeconds
			waitIntervalSeconds := ddosify.IntervalTimeToSeconds(waitInterval)
			log.Printf("Wait interval: %ds", waitIntervalSeconds)
			lc := ddosify.NewLatencyChecker(targetURL, numberOfRuns, waitIntervalSeconds, locations, outputLocationsNumber)
			res, err := lc.RunCommandExec()
			// TODO: Make output prettier
			log.Println(res)
			return err
		},
	}

	addExecFlags(cmd)

	return cmd
}
