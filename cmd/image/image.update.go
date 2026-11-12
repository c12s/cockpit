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
	Long: "Send a digest-pinned OS image update command to a node. The node will pull the image, switch via bootc, reboot, and health-gate the result. " +
		"Without --digest the tag is resolved to the digest it points to now, and that digest is sent.",
	Run: executeImageUpdate,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return utils.ValidateRequiredFlags(cmd, []string{constants.NodeIdFlag, constants.ImageRefFlag})
	},
}

func executeImageUpdate(cmd *cobra.Command, args []string) {
	if digest == "" {
		d, err := resolveDigest(imageRef)
		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}
		digest = d
	}
	accepted, err := requestImageUpdate(imageNodeId, imageRef, digest, rebootPolicy)
	if err != nil {
		fmt.Println("Error sending image update request:", err)
		os.Exit(1)
	}
	if accepted {
		fmt.Printf("Node %s accepted the update to %s@%s.\n", imageNodeId, imageRef, digest)
		fmt.Println("Track progress with: cockpit image status")
	} else {
		// gravity records the desired image before notifying the node
		fmt.Printf("Node %s did not acknowledge the update (offline?). The desired image is recorded; the node applies it when it reconnects.\n", imageNodeId)
	}
}

// resolveDigest pins a tag to the digest it points to now, once, at command
// time: the node, the canary gate and rollback keep working with the digest.
func resolveDigest(ref string) (string, error) {
	cat, err := fetchCatalog()
	if err != nil {
		return "", fmt.Errorf("resolving %s: %w", ref, err)
	}
	d, err := lookupDigest(cat, ref)
	if err != nil {
		return "", err
	}
	fmt.Printf("Resolved %s -> %s\n", ref, d)
	return d, nil
}

// lookupDigest finds repo:tag in the registry listing; a registry host in the
// ref is ignored (the listing is per repo).
func lookupDigest(cat applyCatalog, ref string) (string, error) {
	tag := tagFromRef(ref)
	if tag == "" {
		return "", fmt.Errorf("%s has no tag: give a tag or --digest", ref)
	}
	d, ok := cat.images[repoFromRef(ref)+":"+tag]
	if !ok {
		return "", fmt.Errorf("%s is not in the registry (cockpit image list)", ref)
	}
	return d, nil
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
}
