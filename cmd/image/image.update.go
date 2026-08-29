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

var (
	imageNodeId     string
	imageRef        string
	digest          string
	rebootPolicy    string
	imageUpdateResp model.OSImageUpdateResponse
)

var ImageCmd = &cobra.Command{
	Use:     "image",
	Aliases: []string{"img"},
	Short:   "Manage node OS images (bootc)",
	Long:    "Declare desired OS image versions for nodes and inspect OS status.",
}

var ImageUpdateCmd = &cobra.Command{
	Use:     "update",
	Aliases: []string{"switch"},
	Short:   "Update the OS image of a node",
	Long:    "Send a digest-pinned OS image update command to a node. The node will pull the image, switch via bootc, reboot, and health-gate the result.",
	Run:     executeImageUpdate,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return utils.ValidateRequiredFlags(cmd, []string{constants.NodeIdFlag, constants.ImageRefFlag, constants.DigestFlag})
	},
}

func executeImageUpdate(cmd *cobra.Command, args []string) {
	request := model.OSImageUpdateRequest{
		NodeID:       imageNodeId,
		ImageRef:     imageRef,
		Digest:       digest,
		RebootPolicy: rebootPolicy,
	}
	if rebootPolicy == "" {
		request.RebootPolicy = "immediate"
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
		Response:    &imageUpdateResp,
		Timeout:     30 * time.Second,
	})
	if err != nil {
		fmt.Println("Error sending image update request:", err)
		os.Exit(1)
	}

	if imageUpdateResp.Accepted {
		fmt.Printf("Node %s accepted the update to %s@%s.\n", imageNodeId, imageRef, digest)
		fmt.Println("Track progress with: cockpit image status")
	} else {
		fmt.Printf("Node %s did NOT acknowledge the update command (node offline?).\n", imageNodeId)
		os.Exit(1)
	}
}

func init() {
	ImageCmd.AddCommand(ImageUpdateCmd)
	ImageUpdateCmd.Flags().StringVarP(&imageNodeId, constants.NodeIdFlag, constants.NodeIdShorthandFlag, "", constants.NodeIdDescription)
	ImageUpdateCmd.Flags().StringVarP(&imageRef, constants.ImageRefFlag, "", "", constants.ImageRefDescription)
	ImageUpdateCmd.Flags().StringVarP(&digest, constants.DigestFlag, "", "", constants.DigestDescription)
	ImageUpdateCmd.Flags().StringVarP(&rebootPolicy, constants.RebootPolicyFlag, "", "", constants.RebootPolicyDescription)
	ImageUpdateCmd.MarkFlagRequired(constants.NodeIdFlag)
	ImageUpdateCmd.MarkFlagRequired(constants.ImageRefFlag)
	ImageUpdateCmd.MarkFlagRequired(constants.DigestFlag)
}
