package render

import (
	"fmt"
	"github.com/c12s/cockpit/model"
	"strings"
)

func HandleConfigPlacementsResponse(response *model.ConfigGroupPlacementsResponse) {
	fmt.Println("Config Placements:")
	for _, task := range response.Tasks {
		fmt.Printf("%sTask ID: %s%s\n", Bold, task.ID, Reset)
		fmt.Println(strings.Repeat("-", 45))
		fmt.Printf("  Node: %s\n", task.Node)
		fmt.Printf("  Namespace: %s\n", task.Namespace)
		fmt.Printf("  Status: %s\n", task.Status)
		fmt.Printf("  Accepted At: %s\n", task.AcceptedAt)
		fmt.Printf("  Resolved At: %s\n", task.ResolvedAt)
		fmt.Println(strings.Repeat("-", 45))
	}
	println()
}
