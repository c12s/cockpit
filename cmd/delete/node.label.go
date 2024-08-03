package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/c12s/cockpit/aliases"
	"github.com/c12s/cockpit/clients"
	"github.com/c12s/cockpit/constants"
	"github.com/c12s/cockpit/model"
	"github.com/c12s/cockpit/utils"
	"github.com/spf13/cobra"
)

var (
	nodeId       string
	org          string
	key          string
	nodeResponse model.NodeResponse
)

var DeleteNodeLabelsCmd = &cobra.Command{
	Use:     "label",
	Aliases: aliases.LabelAliases,
	Short:   constants.DeleteNodeLabelsShortDesc,
	Long:    constants.DeleteNodeLabelsLongDesc,
	Run:     executeDeleteLabel,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return utils.ValidateRequiredFlags(cmd, []string{constants.NodeIdFlag, constants.OrganizationFlag, constants.KeyFlag})
	},
}

func executeDeleteLabel(cmd *cobra.Command, args []string) {
	err := sendDeleteLabelRequest()

	if err != nil {
		fmt.Println("Error sending delete node label request:", err)
		os.Exit(1)
	}

	fmt.Printf("Label %s deleted successfully.\n", key)
}

func prepareLabelRequest() model.DeleteLabelInput {
	return model.DeleteLabelInput{
		LabelKey: key,
		NodeID:   nodeId,
		Org:      org,
	}
}

func sendDeleteLabelRequest() error {
	token, err := utils.ReadTokenFromFile()
	if err != nil {
		fmt.Printf("Error reading token: %v\n", err)
		os.Exit(1)
	}

	deleteLabelURL := clients.BuildURL("core", "v1", "DeleteLabel")

	input := prepareLabelRequest()

	return utils.SendHTTPRequest(model.HTTPRequestConfig{
		URL:         deleteLabelURL,
		Method:      "DELETE",
		RequestBody: input,
		Token:       token,
		Response:    &nodeResponse,
		Timeout:     30 * time.Second,
	})
}

func init() {
	DeleteNodeLabelsCmd.Flags().StringVarP(&nodeId, constants.NodeIdFlag, constants.NodeIdShorthandFlag, "", constants.NodeIdDescription)
	DeleteNodeLabelsCmd.Flags().StringVarP(&org, constants.OrganizationFlag, constants.OrganizationShorthandFlag, "", constants.OrganizationDescription)
	DeleteNodeLabelsCmd.Flags().StringVarP(&key, constants.KeyFlag, constants.KeyShorthandFlag, "", constants.LabelKeyDescription)

	DeleteNodeLabelsCmd.MarkFlagRequired(constants.NodeIdFlag)
	DeleteNodeLabelsCmd.MarkFlagRequired(constants.OrganizationFlag)
	DeleteNodeLabelsCmd.MarkFlagRequired(constants.KeyFlag)
}
