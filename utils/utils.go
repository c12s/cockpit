package utils

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/c12s/cockpit/model"
	"golang.org/x/term"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"net/http"
	"strings"
	"syscall"
)

const (
	tokenFilePath = "token.txt"
)

func SendHTTPRequest(config model.HTTPRequestConfig) error {
	var requestBody []byte
	var err error
	if config.RequestBody != nil {
		requestBody, err = json.Marshal(config.RequestBody)
		if err != nil {
			return fmt.Errorf("failed to marshal request body: %v", err)
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), config.Timeout)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, config.Method, config.URL, bytes.NewBuffer(requestBody))
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	if config.Token != "" {
		req.Header.Set("Authorization", "Bearer "+config.Token)
	}
	for key, value := range config.Headers {
		req.Header.Set(key, value)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("request failed with status %s: %s", resp.Status, string(bodyBytes))
	}

	if config.Response != nil {
		if err := json.Unmarshal(bodyBytes, config.Response); err != nil {
			return fmt.Errorf("failed to decode response: %v", err)
		}
	}

	return nil
}

func CreateNodeQuery(query string) ([]model.NodeQuery, error) {
	if query == "" {
		return nil, nil
	}

	parts := strings.Fields(query)
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid query format. Please use 'key operation value'")
	}

	labelKey := parts[0]
	shouldBe := parts[1]
	value := parts[2]

	nodeQuery := model.NodeQuery{
		LabelKey: labelKey,
		ShouldBe: shouldBe,
		Value:    value,
	}

	return []model.NodeQuery{nodeQuery}, nil
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

func PromptForPassword() (string, error) {
	fmt.Print("Enter password: ")
	passwordBytes, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return "", fmt.Errorf("failed to read password: %v", err)
	}
	fmt.Println()
	return string(passwordBytes), nil
}

func ReadTokenFromFile() (string, error) {
	token, err := ioutil.ReadFile(tokenFilePath)
	if err != nil {
		return "", fmt.Errorf("failed to read token file: %w", err)
	}
	return string(token), nil
}
