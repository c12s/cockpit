package render

import (
	"encoding/json"
	"fmt"
	"github.com/c12s/cockpit/model"
	"gopkg.in/yaml.v3"
	"strings"
)

func HandleStandaloneConfigResponse(response *model.StandaloneConfigsResponse, outputFormat string) {
	if outputFormat == "json" {
		jsonData, err := json.MarshalIndent(response, "", "  ")
		if err != nil {
			fmt.Printf("Error converting response to JSON: %v\n", err)
			return
		}
		fmt.Println("Standalone Configurations (JSON):")
		fmt.Println(string(jsonData))
	} else {
		yamlData, err := yaml.Marshal(response)
		if err != nil {
			fmt.Printf("Error converting response to YAML: %v\n", err)
			return
		}
		fmt.Println("Standalone Configurations (YAML):")
		fmt.Println(string(yamlData))
	}
}

func HandleSingleConfigDiffResponse(response *model.SingleConfigDiffResponse, outputFormat string) {
	if outputFormat == "json" {
		jsonData, err := json.MarshalIndent(response, "", "  ")
		if err != nil {
			fmt.Printf("Error converting response to JSON: %v\n", err)
			return
		}
		fmt.Println("Config Diff (JSON):")
		fmt.Println(string(jsonData))
	} else {
		yamlData, err := yaml.Marshal(response)
		if err != nil {
			fmt.Printf("Error converting response to YAML: %v\n", err)
			return
		}
		fmt.Println("Config Diff (YAML):")
		fmt.Println(string(yamlData))
	}
}

func HandleStandaloneConfigPlacementsResponse(response *model.ConfigGroupPlacementsResponse) {
	fmt.Println("Standalone Config Placements:")
	for _, task := range response.Tasks {
		fmt.Printf("%sTask ID: %s%s\n", Bold, task.ID, Reset)
		fmt.Println(strings.Repeat("-", 45))
		fmt.Printf("Node: %s\n", task.Node)
		fmt.Printf("Namespace: %s\n", task.Namespace)
		fmt.Printf("Status: %s\n", task.Status)
		fmt.Printf("Accepted At: %s\n", task.AcceptedAt)
		fmt.Printf("Resolved At: %s\n", task.ResolvedAt)
		fmt.Println(strings.Repeat("-", 45))
	}
}

func DisplayStandaloneResponseAsJSON(response *model.StandaloneConfig) {
	jsonData, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		fmt.Printf("Error converting response to JSON: %v\n", err)
		return
	}
	fmt.Println("Deleted Standalone Config (JSON):")
	fmt.Println(string(jsonData))
	fmt.Println("Standalone configuration deleted successfully!")
}
