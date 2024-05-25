package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/c12s/cockpit/model"
	"gopkg.in/yaml.v3"
	"io/ioutil"
)

const (
	// Path to files
	getConfigFilePathJSON = "./config_group_files/single-config.json"
	getConfigFilePathYAML = "./config_group_files/single-config.yaml"
)

func MarshalAppConfigResponseToYAML(response *model.SingleConfigGroupResponse) ([]byte, error) {
	yamlData, err := yaml.Marshal(response)
	if err != nil {
		return nil, fmt.Errorf("failed to convert to YAML: %v", err)
	}
	return yamlData, nil
}

func MarshalConfigGroupResponseToYAML(response *model.ConfigGroupsResponse) ([]byte, error) {
	var buf bytes.Buffer
	enc := yaml.NewEncoder(&buf)
	enc.SetIndent(2)
	err := enc.Encode(response)
	if err != nil {
		return nil, fmt.Errorf("failed to convert to YAML: %v", err)
	}
	err = enc.Close()
	if err != nil {
		return nil, fmt.Errorf("failed to close encoder: %v", err)
	}
	return buf.Bytes(), nil
}

func MarshalConfigGroupDiffResponseToYAML(response *model.ConfigGroupDiffResponse) ([]byte, error) {
	var buf bytes.Buffer
	enc := yaml.NewEncoder(&buf)
	enc.SetIndent(2)
	err := enc.Encode(response)
	if err != nil {
		return nil, fmt.Errorf("failed to convert to YAML: %v", err)
	}
	err = enc.Close()
	if err != nil {
		return nil, fmt.Errorf("failed to close encoder: %v", err)
	}
	return buf.Bytes(), nil
}

func SaveAppConfigResponseToFiles(response *model.SingleConfigGroupResponse, outputFormat string) error {
	if outputFormat == "json" {
		jsonData, err := json.MarshalIndent(response, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to convert to JSON: %v", err)
		}
		err = ioutil.WriteFile(getConfigFilePathJSON, jsonData, 0644)
		if err != nil {
			return fmt.Errorf("failed to write JSON file: %v", err)
		}
		fmt.Printf("App config saved to %s\n", getConfigFilePathJSON)
	} else {
		yamlData, err := MarshalAppConfigResponseToYAML(response)
		if err != nil {
			return fmt.Errorf("failed to convert to YAML: %v", err)
		}
		err = ioutil.WriteFile(getConfigFilePathYAML, yamlData, 0644)
		if err != nil {
			return fmt.Errorf("failed to write YAML file: %v", err)
		}
		fmt.Printf("App config saved to %s\n", getConfigFilePathYAML)
	}

	return nil
}
