package render

import (
	"encoding/json"
	"fmt"
	"github.com/c12s/cockpit/model"
	"gopkg.in/yaml.v3"
)

func RenderResponseAsTabWriter(data interface{}) {
	switch v := data.(type) {
	case model.ConfigGroup:
		RenderConfigGroupTabWriter(v)
	case model.StandaloneConfig:
		RenderStandaloneConfigTabWriter(v)
	case model.ConfigGroupDiffResponse:
		RenderConfigGroupDiffsTabWriter(v)
	case model.StandaloneConfigDiffResponse:
		RenderStandaloneConfigDiffsTabWriter(v)
	case model.SchemaData:
		RenderSchemaTabWriter(v)
	case []model.SchemaVersion:
		RenderSchemaVersionsTabWriter(v)
	case []model.StandaloneConfig:
		RenderStandaloneConfigsTabWriter(v)
	case []model.ConfigGroup:
		RenderConfigGroupsTabWriter(v)
	case []model.Task:
		RenderTasksTabWriter(v)
	case []model.Node:
		RenderNodesTabWriter(v)
	default:
		fmt.Println("Unsupported data type for tabular rendering")
	}
}

func DisplayResponseAsJSONOrYAML(response interface{}, format, message string) {
	switch format {
	case "json":
		jsonData, err := json.MarshalIndent(response, "", "  ")
		if err != nil {
			fmt.Printf("Error converting response to JSON: %v\n", err)
			return
		}
		fmt.Println(string(jsonData))
		if message != "" {
			fmt.Print(message)
		}
	case "yaml":
		yamlData, err := yaml.Marshal(response)
		if err != nil {
			fmt.Printf("Error converting response to YAML: %v\n", err)
			return
		}
		fmt.Println(string(yamlData))
		if message != "" {
			fmt.Print(message)
		}
	default:
		fmt.Printf("Invalid output format: %v. Supported formats are 'json' and 'yaml'\n", format)
	}
}
