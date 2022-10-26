package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/TwiN/go-color"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/rhsyseng/ddosify-tooling/cli-tool/pkg/ddosify"
	"gopkg.in/yaml.v2"
)

func writeOutputTable(lcol ddosify.LatencyCheckerOutputList) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{color.InWhite("Location"), color.InWhite("Average Latency")})
	if len(lcol.Result) > 0 {
		for i := range lcol.Result {
			if i == 0 {
				t.AppendRow([]interface{}{color.InGreen(lcol.Result[i].Location), color.InGreen(fmt.Sprintf("%f", lcol.Result[i].AvgLatency))})
				continue
			}
			t.AppendRow([]interface{}{color.InWhite(lcol.Result[i].Location), color.InWhite(fmt.Sprintf("%f", lcol.Result[i].AvgLatency))})
		}
	}
	t.SetStyle(table.StyleLight)
	t.Render()
}

func writeOutputJson(lcol ddosify.LatencyCheckerOutputList) {
	o, _ := json.MarshalIndent(lcol, "", "    ")
	fmt.Println(string(o))
}

func writeOutputYaml(lcol ddosify.LatencyCheckerOutputList) {
	o, _ := yaml.Marshal(lcol)
	fmt.Println(string(o))
}
