package utils

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/c12s/cockpit/model"
	"gopkg.in/yaml.v3"
)

func SaveSchemaResponseToYAML(response *model.SchemaResponse, filePath string) error {
	if response.SchemaData.Schema != "" {
		yamlData, err := MarshalSchemaResponse(response)
		if err != nil {
			return fmt.Errorf("failed to convert to YAML: %v", err)
		}

		err = ioutil.WriteFile(filePath, yamlData, 0644)
		if err != nil {
			return fmt.Errorf("failed to write YAML file: %v", err)
		}

		fmt.Printf("Schema saved to %s\n", filePath)
	}
	return nil
}

func SaveVersionResponseToYAML(response *model.SchemaVersionResponse, filePath string) error {
	if len(response.SchemaVersions) != 0 {
		println()
		yamlData, err := MarshalSchemaVersionResponse(response)
		if err != nil {
			return fmt.Errorf("failed to convert to YAML: %v", err)
		}

		err = ioutil.WriteFile(filePath, yamlData, 0644)
		if err != nil {
			return fmt.Errorf("failed to write YAML file: %v", err)
		}

		fmt.Printf("Schema saved to %s\n", filePath)
		return nil
	} else {
		return nil
	}
}

func MarshalSchemaResponse(response *model.SchemaResponse) ([]byte, error) {
	type SchemaData struct {
		Schema       string `yaml:"schema"`
		CreationTime string `yaml:"creationTime"`
	}

	tempResponse := struct {
		Message    string     `yaml:"message"`
		SchemaData SchemaData `yaml:"schemaData"`
	}{
		Message: response.Message,
		SchemaData: SchemaData{
			Schema:       response.SchemaData.Schema,
			CreationTime: response.SchemaData.CreationTime,
		},
	}

	var buf bytes.Buffer
	enc := yaml.NewEncoder(&buf)
	enc.SetIndent(2)
	err := enc.Encode(tempResponse)
	if err != nil {
		return nil, fmt.Errorf("failed to convert to YAML: %v", err)
	}
	err = enc.Close()
	if err != nil {
		return nil, fmt.Errorf("failed to close encoder: %v", err)
	}

	yamlString := buf.String()
	yamlString = removeSchemasBlockScalar(yamlString)
	return []byte(yamlString), nil
}

func MarshalSchemaVersionResponse(response *model.SchemaVersionResponse) ([]byte, error) {
	type SchemaData struct {
		Schema       string `yaml:"schema"`
		CreationTime string `yaml:"creationTime"`
	}

	type SchemaVersion struct {
		SchemaDetails model.SchemaDetails `yaml:"schemaDetails"`
		SchemaData    SchemaData          `yaml:"schemaData"`
	}

	tempResponse := struct {
		Message        string          `yaml:"message"`
		SchemaVersions []SchemaVersion `yaml:"schemaVersions"`
	}{
		Message:        response.Message,
		SchemaVersions: []SchemaVersion{},
	}

	for _, version := range response.SchemaVersions {
		tempResponse.SchemaVersions = append(tempResponse.SchemaVersions, SchemaVersion{
			SchemaDetails: version.SchemaDetails,
			SchemaData: SchemaData{
				Schema:       version.SchemaData.Schema,
				CreationTime: version.SchemaData.CreationTime,
			},
		})
	}

	var buf bytes.Buffer
	enc := yaml.NewEncoder(&buf)
	enc.SetIndent(2)
	err := enc.Encode(tempResponse)
	if err != nil {
		return nil, fmt.Errorf("failed to convert to YAML: %v", err)
	}
	err = enc.Close()
	if err != nil {
		return nil, fmt.Errorf("failed to close encoder: %v", err)
	}

	yamlString := buf.String()
	yamlString = removeSchemasVersionsBlockScalar(yamlString)
	return []byte(yamlString), nil
}

func removeSchemasVersionsBlockScalar(yamlStr string) string {
	lines := strings.Split(yamlStr, "\n")
	for i, line := range lines {
		if strings.Contains(line, "schema: |") {
			lines[i] = strings.Replace(line, "|", "", 1)
		}
	}
	return strings.Join(lines, "\n")
}

func removeSchemasBlockScalar(yamlString string) string {
	return strings.ReplaceAll(yamlString, "schema: |\n", "schema: \n")
}

func ReadSchemaFile(filePath string) (string, error) {
	schema, err := ioutil.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to read schema file: %v", err)
	}
	return string(schema), nil
}
