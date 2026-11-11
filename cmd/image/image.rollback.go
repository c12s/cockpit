package cmd

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/c12s/cockpit/constants"
	"github.com/c12s/cockpit/model"
	"github.com/c12s/cockpit/utils"
	"github.com/spf13/cobra"
)

var imageRollbackOrg string

var ImageRollbackCmd = &cobra.Command{
	Use:     "rollback",
	Aliases: []string{"revert"},
	Short:   "Roll a node back to its previous OS deployment, and keep managing it",
	Long: `Declare the node's previous deployment (the rollback entry bootc keeps on
disk, as reported by the node) as its desired image. The node switches back
through the normal update path (attempt record, greenboot gate, parking of
rejected digests), boots it, and stays managed: afterwards it is healthy on
the previous image and drift is corrected again.

Refused when the previous deployment is the image greenboot rejected (after
an automatic rollback it would boot the broken image again), while an update
is in progress, or when the node reports no previous deployment.`,
	Run: executeImageRollback,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return utils.ValidateRequiredFlags(cmd, []string{constants.NodeIdFlag, constants.OrganizationFlag})
	},
}

func executeImageRollback(cmd *cobra.Command, args []string) {
	nodes, err := orgNodesFor(imageRollbackOrg)
	if err != nil {
		fmt.Println("Error fetching nodes:", err)
		os.Exit(1)
	}
	node, ok := nodes[imageNodeId]
	if !ok {
		fmt.Printf("Node %s is not in org %s.\n", imageNodeId, imageRollbackOrg)
		os.Exit(1)
	}
	ref, digest, err := rollbackTarget(node)
	if err != nil {
		fmt.Printf("Cannot roll back node %s: %v\n", imageNodeId, err)
		os.Exit(1)
	}
	fmt.Printf("Rolling node %s back to %s@%s (its previous deployment).\n", imageNodeId, ref, shortDigest(digest))
	accepted, err := applyImageUpdate(imageNodeId, ref, digest, "")
	if err != nil {
		fmt.Println("Error sending the rollback:", err)
		os.Exit(1)
	}
	if !accepted {
		fmt.Println("The node did not acknowledge (offline?). The desired image is recorded; it switches back when it next pulls.")
		return
	}
	fmt.Println("Accepted. Track it with: cockpit image status --org " + imageRollbackOrg)
}

// rollbackTarget picks the image a rollback declares: the node's previous
// deployment. It refuses the cases that would loop or surprise.
func rollbackTarget(node model.Node) (string, string, error) {
	st := node.OsStatus
	if st == nil {
		return "", "", errors.New("the node sent no OS report yet")
	}
	if st.Phase == "switching" {
		return "", "", errors.New("an update is in progress; wait until it finishes")
	}
	if st.RollbackDigest == "" || st.RollbackImage == "" {
		return "", "", errors.New("the node reports no previous deployment (or runs a star that does not report it); use image update with the digest you want")
	}
	if last := st.LastAttempt; last != nil && last.Result == resultRolledBack && last.Digest == st.RollbackDigest {
		return "", "", fmt.Errorf("the previous deployment %s is the image greenboot rejected; declare a known-good image with image update", shortDigest(st.RollbackDigest))
	}
	ref, _, _ := strings.Cut(st.RollbackImage, "@")
	return ref, st.RollbackDigest, nil
}

func init() {
	ImageCmd.AddCommand(ImageRollbackCmd)
	ImageRollbackCmd.Flags().StringVarP(&imageNodeId, constants.NodeIdFlag, constants.NodeIdShorthandFlag, "", constants.NodeIdDescription)
	ImageRollbackCmd.Flags().StringVarP(&imageRollbackOrg, constants.OrganizationFlag, constants.OrganizationShorthandFlag, "", constants.OrganizationDescription)
	ImageRollbackCmd.MarkFlagRequired(constants.NodeIdFlag)
	ImageRollbackCmd.MarkFlagRequired(constants.OrganizationFlag)
}
