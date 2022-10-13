package main

import (
	"github.com/TwiN/go-color"
	"github.com/rhsyseng/ddosify-tooling/cli-tool/cmd"

	"log"
	"os"

	"github.com/spf13/cobra"
)

func main() {
	command := newCommand()
	if err := command.Execute(); err != nil {
		log.Fatalf(color.InRed("[ERROR]")+"%s", err.Error())
	}
}

func newCommand() *cobra.Command {
	c := &cobra.Command{
		Use:   "ddosify-latencies",
		Short: "ddosify-latencies is the command line interface to work with the ddosify latencies API",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
			os.Exit(1)
		},
	}

	c.AddCommand(cmd.NewExecCommand())

	return c
}
