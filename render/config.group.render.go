package render

import (
	"encoding/json"
	"fmt"
	"github.com/c12s/cockpit/model"
	"gopkg.in/yaml.v3"
	"strings"
)

func HandleSingleConfigGroupResponse(response *model.SingleConfigGroupResponse, outputFormat string) {
	if outputFormat == "json" {
		jsonData, err := json.MarshalIndent(response, "", "  ")
		if err != nil {
			fmt.Printf("Error converting response to JSON: %v\n", err)
			return
		}
		fmt.Println("Config Group (JSON):")
		fmt.Println(string(jsonData))
	} else {
		yamlData, err := yaml.Marshal(response)
		if err != nil {
			fmt.Printf("Error converting response to YAML: %v\n", err)
			return
		}
		fmt.Println("Config group (YAML):")
		fmt.Println(string(yamlData))
	}
}

func HandleConfigGroupResponse(response *model.ConfigGroupsResponse, outputFormat string) {
	if outputFormat == "json" {
		jsonData, err := json.MarshalIndent(response, "", "  ")
		if err != nil {
			fmt.Printf("Error converting response to JSON: %v\n", err)
			return
		}
		fmt.Println("Config Group (JSON):")
		fmt.Println(string(jsonData))
	} else {
		yamlData, err := yaml.Marshal(response)
		if err != nil {
			fmt.Printf("Error converting response to YAML: %v\n", err)
			return
		}
		fmt.Println("Config group (YAML):")
		fmt.Println(string(yamlData))
	}
}

func HandleConfigGroupPlacementsResponse(response *model.ConfigGroupPlacementsResponse) {
	fmt.Println("Config Group Placements:")
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
}

func HandleConfigGroupDiffResponse(response *model.ConfigGroupDiffResponse, outputFormat string) {
	if outputFormat == "json" {
		jsonData, err := json.MarshalIndent(response, "", "  ")
		if err != nil {
			fmt.Printf("Error converting response to JSON: %v\n", err)
			return
		}
		fmt.Println("Config Group Diff (JSON):")
		fmt.Println(string(jsonData))
	} else {
		yamlData, err := yaml.Marshal(response)
		if err != nil {
			fmt.Printf("Error converting response to YAML: %v\n", err)
			return
		}
		fmt.Println("Config Group Diff (YAML):")
		fmt.Println(string(yamlData))
	}
}
