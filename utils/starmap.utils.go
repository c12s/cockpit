package utils

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

func ReadStarmapFile(filePath string) (map[string]interface{}, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}

	ext := strings.ToLower(filepath.Ext(filePath))

	switch ext {
	case ".yaml", ".yml":
		if err := yaml.Unmarshal(data, &result); err != nil {
			return nil, err
		}
	case ".json":
		if err := json.Unmarshal(data, &result); err != nil {
			return nil, err
		}
	default:

		if err := yaml.Unmarshal(data, &result); err != nil {
			if jsonErr := json.Unmarshal(data, &result); jsonErr != nil {
				return nil, err
			}
		}
	}

	return result, nil
}
