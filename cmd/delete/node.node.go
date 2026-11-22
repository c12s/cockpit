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

var deleteNodeResponse model.DeleteNodeResponse

var DeleteNodeCmd = &cobra.Command{
	Use:     "node",
	Aliases: aliases.NodesAliases,
	Short:   constants.DeleteNodeShortDesc,
	Long:    constants.DeleteNodeLongDesc,
	Run:     executeDeleteNode,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return utils.ValidateRequiredFlags(cmd, []string{constants.NodeIdFlag})
	},
}

func executeDeleteNode(cmd *cobra.Command, args []string) {
	if err := sendDeleteNodeRequest(); err != nil {
		fmt.Println("Error sending delete node request:", err)
		os.Exit(1)
	}

	fmt.Printf("Node %s deleted successfully.\n", nodeId)
}

func sendDeleteNodeRequest() error {
	token, err := utils.ReadTokenFromFile()
	if err != nil {
		return fmt.Errorf("error reading token: %v", err)
	}

	deleteNodeURL := clients.BuildURL("core", "v1", "DeleteNode")

	return utils.SendHTTPRequest(model.HTTPRequestConfig{
		URL:         deleteNodeURL,
		Method:      "DELETE",
		RequestBody: model.DeleteNodeInput{NodeID: nodeId},
		Token:       token,
		Response:    &deleteNodeResponse,
		Timeout:     30 * time.Second,
	})
}

func init() {
	DeleteNodeCmd.Flags().StringVarP(&nodeId, constants.NodeIdFlag, constants.NodeIdShorthandFlag, "", constants.NodeIdDescription)
	DeleteNodeCmd.MarkFlagRequired(constants.NodeIdFlag)
}
