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
	imageNodeId  string
	imageRef     string
	digest       string
	rebootPolicy string
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
	accepted, err := requestImageUpdate(imageNodeId, imageRef, digest, rebootPolicy)
	if err != nil {
		fmt.Println("Error sending image update request:", err)
		os.Exit(1)
	}
	if accepted {
		fmt.Printf("Node %s accepted the update to %s@%s.\n", imageNodeId, imageRef, digest)
		fmt.Println("Track progress with: cockpit image status")
	} else {
		fmt.Printf("Node %s did NOT acknowledge the update command (node offline?).\n", imageNodeId)
		os.Exit(1)
	}
}

// requestImageUpdate sends a digest-pinned update to one node and reports
// whether the node acknowledged it.
func requestImageUpdate(nodeID, ref, digest, policy string) (bool, error) {
	if policy == "" {
		policy = "immediate"
	}
	token, err := utils.ReadTokenFromFile()
	if err != nil {
		return false, fmt.Errorf("reading token: %w", err)
	}
	var resp model.OSImageUpdateResponse
	err = utils.SendHTTPRequest(model.HTTPRequestConfig{
		URL:    clients.BuildURL("core", "v1", "DisseminateOSImage"),
		Method: "POST",
		Token:  token,
		RequestBody: model.OSImageUpdateRequest{
			NodeID:       nodeID,
			ImageRef:     ref,
			Digest:       digest,
			RebootPolicy: policy,
		},
		Response: &resp,
		Timeout:  30 * time.Second,
	})
	if err != nil {
		return false, err
	}
	return resp.Accepted, nil
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
