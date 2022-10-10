package cmd

import (
	"github.com/rhsyseng/ddosify-tooling/cli-tool/pkg/ddosify"
	"github.com/spf13/cobra"
)

func addRunFlags(cmd *cobra.Command) {

	flags := cmd.Flags()
	flags.String("target-url", "https://google.com", "The target url.")
}

func NewStatusCommand() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "status",
		Short: "Run the status command",
		RunE: func(cmd *cobra.Command, args []string) error {
			return ddosify.RunCommandStatus()
		},
	}

	addRunFlags(cmd)

	return cmd
}
