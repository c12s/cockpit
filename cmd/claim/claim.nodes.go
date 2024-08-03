package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/c12s/cockpit/aliases"
	"github.com/c12s/cockpit/clients"
	"github.com/c12s/cockpit/constants"
	"github.com/c12s/cockpit/model"
	"github.com/c12s/cockpit/render"
	"github.com/c12s/cockpit/utils"
	"github.com/cheggaaa/pb/v3"
	"github.com/spf13/cobra"
)

var (
	org               string
	query             string
	claimNodeResponse model.ClaimNodesResponse
)

var ClaimNodesCmd = &cobra.Command{
	Use:     "nodes",
	Aliases: aliases.ClaimAliases,
	Short:   constants.ClaimNodesShortDesc,
	Long:    constants.ClaimNodesLongDesc,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return utils.ValidateRequiredFlags(cmd, []string{constants.OrganizationFlag, constants.QueryFlag})
	},
	Run: executeClaimNodes,
}

func executeClaimNodes(cmd *cobra.Command, args []string) {
	requestBody, err := prepareClaimNodesRequest()
	if err != nil {
		fmt.Println("Error preparing request:", err)
		os.Exit(1)
	}

	bar := pb.StartNew(100)
	bar.SetTemplate(pb.Full)
	bar.SetWidth(50)

	go func() {
		for {
			time.Sleep(50 * time.Millisecond)
			bar.Increment()
		}
	}()

	if err := sendClaimNodeRequest(requestBody, bar); err != nil {
		bar.Finish()
		fmt.Println("Error claiming nodes:", err)
		os.Exit(1)
	}

	bar.Finish()
	render.RenderResponseAsTabWriter(claimNodeResponse.Nodes)
	fmt.Println("These nodes were successfully claimed!")
}

func prepareClaimNodesRequest() (interface{}, error) {
	request := model.ClaimNodesRequest{
		Org: org,
	}

	nodeQueries, err := utils.CreateNodeQuery(query)
	if err != nil {
		return nil, err
	}
	request.Query = nodeQueries

	return request, nil
}

func sendClaimNodeRequest(requestBody interface{}, bar *pb.ProgressBar) error {
	token, err := utils.ReadTokenFromFile()
	if err != nil {
		return fmt.Errorf("error reading token: %v", err)
	}

	url := clients.BuildURL("core", "v1", "ClaimOwnership")

	err = utils.SendHTTPRequestWithProgress(model.HTTPRequestConfig{
		URL:         url,
		Method:      "PATCH",
		Token:       token,
		RequestBody: requestBody,
		Response:    &claimNodeResponse,
		Timeout:     30 * time.Second,
	}, bar)

	if err != nil {
		return err
	}

	bar.SetCurrent(100)
	return nil
}

func init() {
	ClaimNodesCmd.Flags().StringVarP(&org, constants.OrganizationFlag, constants.OrganizationShorthandFlag, "", constants.OrganizationDescription)
	ClaimNodesCmd.Flags().StringVarP(&query, constants.QueryFlag, constants.QueryFlagShorthandFlag, "", constants.NodeQueryRequiredDescription)

	ClaimNodesCmd.MarkFlagRequired(constants.OrganizationFlag)
	ClaimNodesCmd.MarkFlagRequired(constants.QueryFlag)
}
