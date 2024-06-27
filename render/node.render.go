package render

import (
	"fmt"
	"github.com/c12s/cockpit/model"
	"os"
	"strconv"
	"strings"
	"text/tabwriter"
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

func RenderNodesTabWriter(nodes []model.Node) {
	if len(nodes) == 0 {
		fmt.Println("No nodes were found.")
		return
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', tabwriter.TabIndent)
	defer w.Flush()

	fmt.Fprintf(w, "Node ID\tCPU Cores\tAvg Clock Speed (MHz)\tAvg Cache (KB)\tMemory (GB)\tDisk Total (GB)\tDisk Free (GB)\n")

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

		fmt.Fprintf(w, "%s\t%.2f\t%.2f\t%.2f\t%.2f\t%.2f\t%.2f\n", node.ID, cpuCores, avgMhz, avgCache, memoryTotalGB, diskTotalGB, diskFreeGB)
	}
}
