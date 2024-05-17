package model

import (
	"fmt"
	"strings"
)

func RenderNodes(nodes []Node) {
	for _, node := range nodes {
		fmt.Printf("Node ID: %s\n", node.ID)
		fmt.Println(strings.Repeat("-", 45))
		for _, label := range node.Labels {
			fmt.Printf("  - %s: %s\n", label.Key, label.Value)
		}
		fmt.Println(strings.Repeat("-", 45))
	}
}
