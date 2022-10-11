package cmd

import (
	"github.com/rhsyseng/ddosify-tooling/cli-tool/pkg/ddosify"
	"github.com/spf13/cobra"
)

var (
	targetURL string
)

func addRunFlags(cmd *cobra.Command) {

	flags := cmd.Flags()
	flags.StringVar(&targetURL, "target-url", "https://google.com", "The target url.")
}

func NewStatusCommand() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "status",
		Short: "Run the status command",
		RunE: func(cmd *cobra.Command, args []string) error {
			lc := ddosify.NewLatencyChecker(targetURL)
			return lc.RunCommandStatus()
		},
	}

	addRunFlags(cmd)

	return cmd
}
