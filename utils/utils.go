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

func SaveConfigResponseToFile(response interface{}, filePath string) error {
	if strings.HasSuffix(filePath, ".json") {
		jsonData, err := json.MarshalIndent(response, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to convert to JSON: %v", err)
		}
		err = ioutil.WriteFile(filePath, jsonData, 0644)
		if err != nil {
			return fmt.Errorf("failed to write JSON file: %v", err)
		}
		fmt.Printf("Response saved to %s\n", filePath)
	} else if strings.HasSuffix(filePath, ".yaml") {
		yamlData, err := yaml.Marshal(response)
		if err != nil {
			return fmt.Errorf("failed to convert to YAML: %v", err)
		}
		err = ioutil.WriteFile(filePath, yamlData, 0644)
		if err != nil {
			return fmt.Errorf("failed to write YAML file: %v", err)
		}
		fmt.Printf("Response saved to %s\n", filePath)
	} else {
		return fmt.Errorf("unsupported file extension")
	}
	return nil
}
