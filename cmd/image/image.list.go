package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/c12s/cockpit/clients"
	"github.com/c12s/cockpit/constants"
	"github.com/c12s/cockpit/model"
	"github.com/c12s/cockpit/utils"
	"github.com/spf13/cobra"
)

// ImageListCmd lists the images and qcow2 artifacts available through the
// control plane. It goes through the gateway's ListImages route, so it uses
// the same token and policy path as the other image commands.
var ImageListCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List images available in the registry and artifact server",
	Long: "List the images in the registry and the qcow2 artifacts in the artifact " +
		"server. The command is read-only and goes through the gateway's " +
		"ListImages route, so it uses the same token as the other commands.",
	Run: executeImageList,
}

var imageListOutput string

func executeImageList(cmd *cobra.Command, args []string) {
	token, err := utils.ReadTokenFromFile()
	if err != nil {
		fmt.Printf("Error reading token: %v\n", err)
		os.Exit(1)
	}

	var resp model.ImageListResponse
	err = utils.SendHTTPRequest(model.HTTPRequestConfig{
		URL:      clients.BuildURL("core", "v1", "ListImages"),
		Method:   "GET",
		Token:    token,
		Response: &resp,
		Timeout:  30 * time.Second,
	})
	if err != nil {
		fmt.Printf("Error listing images: %v\n", err)
		os.Exit(1)
	}

	if imageListOutput == "json" {
		out, _ := json.MarshalIndent(resp, "", "  ")
		fmt.Println(string(out))
		return
	}

	fmt.Println("Images")
	if len(resp.Images) == 0 {
		fmt.Println("  (none)")
	} else {
		for _, img := range resp.Images {
			fmt.Printf("  %s\n", img)
		}
	}

	fmt.Println("\nArtifacts")
	if len(resp.Artifacts) == 0 {
		fmt.Println("  (none)")
	} else {
		for _, a := range resp.Artifacts {
			fmt.Printf("  %s  %s\n", a.Name, a.Checksum)
		}
	}
}

func init() {
	ImageListCmd.Flags().StringVarP(&imageListOutput, constants.OutputFlag, constants.OutputShorthandFlag, "table", "Output format: table or json")
	ImageCmd.AddCommand(ImageListCmd)
}
