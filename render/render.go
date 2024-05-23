package render

import (
	"fmt"
	"github.com/c12s/cockpit/model"
	"strings"
)

const (
	Bold  = "\033[1m"
	Reset = "\033[0m"
)

func HandleSchemaVersionResponse(response *model.SchemaVersionResponse) {
	println()
	fmt.Println("Message:", response.Message)
	for _, version := range response.SchemaVersions {
		fmt.Println("Schema Name:", version.SchemaDetails.SchemaName)
		fmt.Println("Version:", version.SchemaDetails.Version)
		fmt.Println("Organization:", version.SchemaDetails.Organization)
		fmt.Println("Schema Data:")
		fmt.Println(version.SchemaData.Schema)
		fmt.Println("Creation Time:", version.SchemaData.CreationTime)
		fmt.Println()
	}
}

func HandleSchemaResponse(response *model.SchemaResponse) {
	println()
	fmt.Println("Message:", response.Message)
	if response.SchemaData.Schema != "" {
		fmt.Println("Schema Data:")
		fmt.Println(response.SchemaData.Schema)
		fmt.Println("Creation Time:", response.SchemaData.CreationTime)
	}
	println()
}

func RenderNodes(nodes []model.Node) {
	if len(nodes) == 0 {
		fmt.Println("No nodes were found.")
	} else {
		for _, node := range nodes {
			fmt.Printf("%sNode ID: %s%s\n", Bold, Bold, node.ID)
			fmt.Println(strings.Repeat("-", 45))
			for _, label := range node.Labels {
				fmt.Printf("  - %s%s: %s%s\n", Reset, label.Key, Reset, label.Value)
			}
			fmt.Println(strings.Repeat("-", 45))
		}
	}
}

func RenderNode(node model.Node) {
	fmt.Printf("%sNode ID: %s%s\n", Bold, Bold, node.ID)
	fmt.Println(strings.Repeat("-", 45))
	for _, label := range node.Labels {
		fmt.Printf("  - %s%s: %s%s\n", Reset, label.Key, Reset, label.Value)
	}
	fmt.Println(strings.Repeat("-", 45))
}
