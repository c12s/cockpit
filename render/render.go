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

func RenderNodes(nodes []model.Node) {
	for _, node := range nodes {
		fmt.Printf("%sNode ID: %s%s\n", Bold, Bold, node.ID)
		fmt.Println(strings.Repeat("-", 45))
		for _, label := range node.Labels {
			fmt.Printf("  - %s%s: %s%s\n", Reset, label.Key, Reset, label.Value)
		}
		fmt.Println(strings.Repeat("-", 45))
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
