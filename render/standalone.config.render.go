package render

import (
	"encoding/json"
	"fmt"
	"github.com/c12s/cockpit/model"
	"gopkg.in/yaml.v3"
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
