package render

import (
	"encoding/json"
	"fmt"
	"github.com/c12s/cockpit/model"
	"gopkg.in/yaml.v3"
)

func RenderResponseToYAMLOrJSON(response interface{}, outputFormat string) {
	println()
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
	println()
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
	println()
}
