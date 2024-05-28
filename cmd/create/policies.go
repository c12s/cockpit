package cmd

import (
	"fmt"
	"github.com/c12s/cockpit/clients"
	"github.com/c12s/cockpit/model"
	"github.com/c12s/cockpit/utils"
	"log"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

const (
	createPoliciesShortDesc = "Create policies from YAML or JSON file"
	createPoliciesLongDesc  = "This command  is for creating security policies based on the input file.\n\n" +
		"Example:\n" +
		"cockpit create policies --path 'path to yaml or json file'"
)

var CreatePoliciesCmd = &cobra.Command{
	Use:   "policies",
	Short: createPoliciesShortDesc,
	Long:  createPoliciesLongDesc,
	Run:   executeCreatePolicies,
}

func executeCreatePolicies(cmd *cobra.Command, args []string) {
	requestBody, err := preparePoliciesRequestBody()
	if err != nil {
		log.Fatal(err)
	}

	if err := sendCreatePoliciesRequest(requestBody); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Policies created successfully")
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
