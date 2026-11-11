package render

import (
	"fmt"
	"github.com/c12s/cockpit/model"
	"os"
	"strconv"
	"strings"
	"text/tabwriter"
	"time"
)

func RenderNodes(nodes []model.Node) {
	if len(nodes) == 0 {
		fmt.Println("No nodes were found.")
	} else {
		for _, node := range nodes {
			fmt.Printf("Node ID: %s\n", node.ID)
			fmt.Println(strings.Repeat("-", 45))
			if node.Liveness != "" {
				fmt.Printf("  liveness: %s (%s)\n", node.Liveness, FormatLastSeen(node.LastSeen))
			}
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
	if node.Liveness != "" {
		fmt.Printf("  liveness: %s (%s)\n", node.Liveness, FormatLastSeen(node.LastSeen))
	}
	for _, label := range node.Labels {
		fmt.Printf("  - %s: %s\n", label.Key, label.Value)
	}
	fmt.Println(strings.Repeat("-", 45))
}

func RenderNodesTabWriter(nodes []model.Node) {
	if len(nodes) == 0 {
		fmt.Println("No nodes were found.")
		return
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', tabwriter.TabIndent)
	defer w.Flush()

	fmt.Fprintf(w, "Node ID\tLiveness\tCPU Cores\tAvg Clock Speed (MHz)\tAvg Cache (KB)\tMemory (GB)\tDisk Total (GB)\tDisk Free (GB)\n")

	for _, node := range nodes {
		var cpuCores, totalMhz, totalCache, memoryTotalGB, diskTotalGB, diskFreeGB float64
		var coreCount int

		for _, label := range node.Labels {
			value, ok := label.Value.(string)
			if !ok {
				continue
			}
			switch label.Key {
			case "cpu-cores":
				cpuCores, _ = strconv.ParseFloat(value, 64)
			case "memory-totalGB":
				memoryTotalGB, _ = strconv.ParseFloat(value, 64)
			case "disk-totalGB":
				diskTotalGB, _ = strconv.ParseFloat(value, 64)
			case "disk-freeGB":
				diskFreeGB, _ = strconv.ParseFloat(value, 64)
			}

			if strings.Contains(label.Key, "mhz") {
				mhz, _ := strconv.ParseFloat(value, 64)
				totalMhz += mhz
				coreCount++
			}

			if strings.Contains(label.Key, "cacheKB") {
				cache, _ := strconv.ParseFloat(value, 64)
				totalCache += cache
			}
		}

		avgMhz := totalMhz / float64(coreCount)
		avgCache := totalCache / float64(coreCount)

		liveness := node.Liveness
		if liveness == "" {
			liveness = "-"
		}
		fmt.Fprintf(w, "%s\t%s\t%.2f\t%.2f\t%.2f\t%.2f\t%.2f\t%.2f\n", node.ID, liveness, cpuCores, avgMhz, avgCache, memoryTotalGB, diskTotalGB, diskFreeGB)
	}
}

// FormatLastSeen renders a unix timestamp as "just now", a relative age, or
// "-" when the node has never been seen.
func FormatLastSeen(lastSeen model.FlexInt64) string {
	if lastSeen <= 0 {
		return "never seen"
	}
	age := time.Since(time.Unix(int64(lastSeen), 0))
	if age < time.Second {
		return "just now"
	}
	return age.Round(time.Second).String() + " ago"
}
