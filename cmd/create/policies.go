package cmd

import (
	"fmt"
	"github.com/c12s/cockpit/clients"
	"github.com/c12s/cockpit/model"
	"github.com/c12s/cockpit/utils"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

const (
	createPoliciesShortDesc = "Create policies from YAML or JSON file"
	createPoliciesLongDesc  = `This command is for creating security policies based on the input file.
Policies are used to define and enforce security rules within the organization. The input file can be in YAML or JSON format, specifying the policy details.

Example:
cockpit create policies --path 'path to yaml or json file'`
)

var CreatePoliciesCmd = &cobra.Command{
	Use:     "policies",
	Aliases: []string{"policie", "policiess", "pol"},
	Short:   createPoliciesShortDesc,
	Long:    createPoliciesLongDesc,
	Run:     executeCreatePolicies,
}

func executeCreatePolicies(cmd *cobra.Command, args []string) {
	requestBody, err := preparePoliciesRequestBody()
	if err != nil {
		fmt.Println("Error preparing request:", err)
		os.Exit(1)
	}

	if err := sendCreatePoliciesRequest(requestBody); err != nil {
		fmt.Println("Error sending policies request:", err)
		os.Exit(1)
	}

	fmt.Println("Policies created successfully")
	fmt.Println()
}

func preparePoliciesRequestBody() (model.PoliciesRequest, error) {
	var requestBody model.PoliciesRequest
	var err error
	if strings.HasSuffix(filePath, ".yaml") {
		err = utils.ReadYAML(filePath, &requestBody)
	} else if strings.HasSuffix(filePath, ".json") {
		err = utils.ReadJSON(filePath, &requestBody)
	} else {
		return requestBody, fmt.Errorf("invalid file format. Please provide a YAML or JSON file")
	}

	if err != nil {
		return requestBody, fmt.Errorf("failed to read input file: %v", err)
	}

	return requestBody, nil
}

func sendCreatePoliciesRequest(requestBody model.PoliciesRequest) error {
	config, err := preparePoliciesRequestConfig(requestBody)
	if err != nil {
		return fmt.Errorf("error creating request config: %v", err)
	}

	err = utils.SendHTTPRequest(config)
	if err != nil {
		return fmt.Errorf("failed to send request: %v", err)
	}

	return nil
}

func preparePoliciesRequestConfig(requestBody model.PoliciesRequest) (model.HTTPRequestConfig, error) {
	token, err := utils.ReadTokenFromFile()
	if err != nil {
		return model.HTTPRequestConfig{}, fmt.Errorf("error reading token: %v", err)
	}

	url := clients.BuildURL("core", "v1", "CreatePolicy")

	return model.HTTPRequestConfig{
		URL:         url,
		Method:      "POST",
		Token:       token,
		RequestBody: requestBody,
		Timeout:     10 * time.Second,
	}, nil
}

func init() {
	CreatePoliciesCmd.Flags().StringVarP(&filePath, filePathFlag, "p", "", filePathDescription)
	CreatePoliciesCmd.MarkFlagRequired(filePathFlag)
}
