package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/c12s/cockpit/clients"
	"github.com/c12s/cockpit/constants"
	"github.com/c12s/cockpit/model"
	"github.com/c12s/cockpit/utils"
	"github.com/spf13/cobra"
)

var imageRollbackResp model.OSImageUpdateResponse

var ImageRollbackCmd = &cobra.Command{
	Use:     "rollback",
	Aliases: []string{"revert"},
	Short:   "Roll back the OS image of a node",
	Long:    "Order a node to boot into its previous OS deployment via bootc rollback. The node clears its desired image and reboots into the last known-good OS.",
	Run:     executeImageRollback,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return utils.ValidateRequiredFlags(cmd, []string{constants.NodeIdFlag})
	},
}

func executeImageRollback(cmd *cobra.Command, args []string) {
	request := model.OSImageUpdateRequest{
		NodeID:   imageNodeId,
		Rollback: true,
	}

	url := clients.BuildURL("core", "v1", "DisseminateOSImage")
	token, err := utils.ReadTokenFromFile()
	if err != nil {
		fmt.Println("Error reading token:", err)
		os.Exit(1)
	}

	err = utils.SendHTTPRequest(model.HTTPRequestConfig{
		URL:         url,
		Method:      "POST",
		Token:       token,
		RequestBody: request,
		Response:    &imageRollbackResp,
		Timeout:     30 * time.Second,
	})
	if err != nil {
		fmt.Println("Error sending image rollback request:", err)
		os.Exit(1)
	}

	if imageRollbackResp.Accepted {
		fmt.Printf("Node %s accepted the rollback command.\n", imageNodeId)
		fmt.Println("Track progress with: cockpit image status")
	} else {
		fmt.Printf("Node %s did NOT acknowledge the rollback command (node offline?).\n", imageNodeId)
		os.Exit(1)
	}
}

func init() {
	ImageCmd.AddCommand(ImageRollbackCmd)
	ImageRollbackCmd.Flags().StringVarP(&imageNodeId, constants.NodeIdFlag, constants.NodeIdShorthandFlag, "", constants.NodeIdDescription)
	ImageRollbackCmd.MarkFlagRequired(constants.NodeIdFlag)
}
