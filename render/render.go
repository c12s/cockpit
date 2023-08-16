package render

import (
	"fmt"
	"github.com/c12s/cockpit/model"
	"github.com/c12s/magnetar/pkg/api"
	"github.com/fatih/color"
	"github.com/jedib0t/go-pretty/v6/table"
	"os"
)

func Error(err error) {
	color.Red(err.Error())
}

func Config(config model.ConfigGroup) {
	fmt.Printf("Config name: %s\n", config.Name)
	fmt.Printf("Organization: %s\n", config.OrgId)
	fmt.Printf("Version: v%d\n", config.Version)
	fmt.Printf("Configs: \n")
	for key, value := range config.Configs {
		fmt.Printf("\t%s = %s\n", key, value)
	}
}

func Node(node *api.NodeStringified) {
	Nodes([]*api.NodeStringified{node})
}

func Nodes(nodes []*api.NodeStringified) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"Node ID", "Labels"})
	for _, node := range nodes {
		t.AppendRow(table.Row{node.Id, formatLabels(node.Labels)})
		t.AppendSeparator()
	}
	t.Render()
}

func formatLabels(labels []*api.LabelStringified) (formatted string) {
	for i, label := range labels {
		if i > 0 {
			formatted += "\n"
		}
		formatted += fmt.Sprintf("%s=%s", label.Key, label.Value)
	}
	return
}
