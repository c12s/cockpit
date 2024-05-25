package utils

import (
	"bytes"
	"fmt"
	"github.com/c12s/cockpit/model"
	"gopkg.in/yaml.v3"
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
