package render

import (
	"encoding/json"
	"fmt"
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

func DisplayResponse(response interface{}, format, message string) {
	switch format {
	case "json":
		jsonData, err := json.MarshalIndent(response, "", "  ")
		if err != nil {
			fmt.Printf("Error converting response to JSON: %v\n", err)
			return
		}
		fmt.Println(message)
		fmt.Println(string(jsonData))
	case "yaml":
		yamlData, err := yaml.Marshal(response)
		if err != nil {
			fmt.Printf("Error converting response to YAML: %v\n", err)
			return
		}
		fmt.Println(message)
		fmt.Println(string(yamlData))
	default:
		fmt.Printf("Invalid output format: %v. Supported formats are 'json' and 'yaml'\n", format)
	}
	fmt.Println("Deleted successfully!")
}
