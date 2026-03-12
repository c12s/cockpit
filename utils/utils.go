package utils

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"syscall"

	"github.com/c12s/cockpit/model"
	"github.com/cheggaaa/pb/v3"
	"github.com/spf13/cobra"
	"golang.org/x/term"
	"gopkg.in/yaml.v3"
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

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("request failed with status %s", string(bodyBytes))
	}

	if config.Response != nil {
		if err := json.Unmarshal(bodyBytes, config.Response); err != nil {
			return fmt.Errorf("failed to decode response: %v", err)
		}
	}

	return nil
}

func SendHTTPRequestWithProgress(config model.HTTPRequestConfig, bar *pb.ProgressBar) error {
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

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("request failed with status %s", string(bodyBytes))
	}

	if config.Response != nil {
		if err := json.Unmarshal(bodyBytes, config.Response); err != nil {
			return fmt.Errorf("failed to decode response: %v", err)
		}
	}

	bar.SetCurrent(100)
	bar.Finish()
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

func SaveTokenToFile(token string) error {
	tokenFilePath := tokenFilePath
	return ioutil.WriteFile(tokenFilePath, []byte(token), 0600)
}

func SaveYAMLOrJSONResponseToFile(response interface{}, filePath string) error {
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

func ReadYAML(filePath string, out interface{}) error {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(data, out)
	if err != nil {
		return err
	}

	return nil
}

func ReadJSON(filePath string, out interface{}) error {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}

	err = json.Unmarshal(data, out)
	if err != nil {
		return err
	}

	return nil
}

func PrepareRequestBodyFromYAMLOrJSON(path string) (map[string]interface{}, error) {
	var configData map[string]interface{}

	fileContent, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %v", err)
	}

	if strings.HasSuffix(path, ".yaml") {
		err = yaml.Unmarshal(fileContent, &configData)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal YAML: %v", err)
		}
	} else if strings.HasSuffix(path, ".json") {
		err = json.Unmarshal(fileContent, &configData)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal JSON: %v", err)
		}
	} else {
		return nil, fmt.Errorf("unsupported file format")
	}

	return configData, nil
}

func LoadEnvFile(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.TrimSpace(line) == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			return fmt.Errorf("invalid line in .env file: %s", line)
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		err = os.Setenv(key, value)
		if err != nil {
			return err
		}
	}

	return scanner.Err()
}

func StringToFloat(s string) float64 {
	value, _ := strconv.ParseFloat(s, 64)
	return value
}

func ValidateRequiredFlags(cmd *cobra.Command, requiredFlags []string) error {
	for _, flag := range requiredFlags {
		value, err := cmd.Flags().GetString(flag)
		if err != nil {
			return err
		}
		if strings.TrimSpace(value) == "" {
			return errors.New("required flag --" + flag + " cannot be empty")
		}
	}
	return nil
}

func ValidateRequiredFlagsAny(cmd *cobra.Command, requiredFlags []string) error {
	for _, flagName := range requiredFlags {
		flag := cmd.Flags().Lookup(flagName)
		if flag == nil {
			return fmt.Errorf("flag --%s not defined", flagName)
		}
		if !flag.Changed {
			return fmt.Errorf("required flag --%s not set", flagName)
		}
		if flag.Value.Type() == "string" && strings.TrimSpace(flag.Value.String()) == "" {
			return fmt.Errorf("required flag --%s cannot be empty", flagName)
		}
	}
	return nil
}
