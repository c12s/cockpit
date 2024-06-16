package render

import (
	"fmt"
	"github.com/c12s/cockpit/model"
	"strings"
)

func RenderNodes(nodes []model.Node) {
	if len(nodes) == 0 {
		fmt.Println("No nodes were found.")
	} else {
		for _, node := range nodes {
			fmt.Printf("Node ID: %s\n", node.ID)
			fmt.Println(strings.Repeat("-", 45))
			for _, label := range node.Labels {
				fmt.Printf("  - %s: %s\n", label.Key, label.Value)
			}
			fmt.Println(strings.Repeat("-", 45))
		}
	}
}

func RenderNode(node model.Node) {
	fmt.Printf("Node ID: %s\n", node.ID)
	fmt.Println(strings.Repeat("-", 45))
	for _, label := range node.Labels {
		fmt.Printf("  - %s: %s\n", label.Key, label.Value)
	}
	fmt.Println(strings.Repeat("-", 45))
}
