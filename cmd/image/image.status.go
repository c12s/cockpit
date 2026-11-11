package cmd

import (
	"fmt"
	"os"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/c12s/cockpit/clients"
	"github.com/c12s/cockpit/constants"
	"github.com/c12s/cockpit/model"
	"github.com/c12s/cockpit/utils"
	"github.com/spf13/cobra"
)

var (
	imageStatusWide  bool
	imageStatusOrg   string
	imageStatusNodes model.NodesResponse
)

var ImageStatusCmd = &cobra.Command{
	Use:     "status",
	Aliases: []string{"st"},
	Short:   "Show OS status of allocated nodes",
	Long:    "Display the current OS image, digest, and update phase (healthy / switching / rolled_back) for all nodes of an organization.",
	Run:     executeImageStatus,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return utils.ValidateRequiredFlags(cmd, []string{constants.OrganizationFlag})
	},
}

func executeImageStatus(cmd *cobra.Command, args []string) {
	url := clients.BuildURL("core", "v1", "ListOrgOwnedNodes")
	token, err := utils.ReadTokenFromFile()
	if err != nil {
		fmt.Println("Error reading token:", err)
		os.Exit(1)
	}

	err = utils.SendHTTPRequest(model.HTTPRequestConfig{
		URL:         url,
		Method:      "GET",
		Token:       token,
		RequestBody: model.ClaimNodesRequest{Org: imageStatusOrg},
		Response:    &imageStatusNodes,
		Timeout:     30 * time.Second,
	})
	if err != nil {
		fmt.Println("Error fetching nodes:", err)
		os.Exit(1)
	}

	fmt.Print(renderStatusTable(imageStatusNodes.Nodes, imageTags(), imageStatusWide))
}

// imageTags maps each registry digest to its tags, so a digest can be shown
// as the version it is. A failed lookup only leaves the tags out.
func imageTags() map[string][]string {
	var resp model.ImageListResponse
	if err := gatewayGet("ListImages", nil, &resp); err != nil {
		return nil
	}
	tags := map[string][]string{}
	for _, img := range resp.Images {
		tags[img.Digest] = append(tags[img.Digest], tagFromRef(img.Ref))
	}
	return tags
}

// describeImage shows a digest as "v4.1 (sha256:b4e5e68be4ec)": the tag(s)
// pointing at it plus the digest, short unless wide.
func describeImage(digest string, tags map[string][]string, wide bool) string {
	if digest == "" {
		return "-"
	}
	d := digest
	if !wide && len(d) > len("sha256:")+12 {
		d = d[:len("sha256:")+12]
	}
	if t := tags[digest]; len(t) > 0 {
		return strings.Join(t, ",") + " (" + d + ")"
	}
	return d
}

// renderStatusTable lays the nodes out with tabwriter, so columns line up
// whatever the value lengths are.
func renderStatusTable(nodes []model.Node, tags map[string][]string, wide bool) string {
	var b strings.Builder
	w := tabwriter.NewWriter(&b, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "NODE ID\tPHASE\tVERSION\tBOOTED\tDESIRED\tLAST ATTEMPT")
	for _, node := range nodes {
		phase, version, booted, attempt := "-", "-", "-", "-"
		// desired = declared in the control plane; "-" means unmanaged
		desired := "-"
		if node.DesiredOs != nil && node.DesiredOs.Digest != "" {
			desired = describeImage(node.DesiredOs.Digest, tags, wide)
		}
		if st := node.OsStatus; st != nil {
			if st.Phase != "" {
				phase = st.Phase
			}
			if st.BootedVersion != "" {
				version = st.BootedVersion
			}
			booted = describeImage(st.BootedDigest, tags, wide)
			if st.LastAttempt != nil {
				attempt = fmt.Sprintf("%s (%s)", st.LastAttempt.Result, time.Unix(int64(st.LastAttempt.Timestamp), 0).Format("15:04:05"))
			}
		}
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%s\n", node.ID, phase, version, booted, desired, attempt)
	}
	w.Flush()
	return b.String()
}

func shortDigest(digest string) string {
	if len(digest) > 19 {
		return digest[:19] + "..."
	}
	return digest
}

func init() {
	ImageCmd.AddCommand(ImageStatusCmd)
	ImageStatusCmd.Flags().StringVarP(&imageStatusOrg, constants.OrganizationFlag, constants.OrganizationShorthandFlag, "", constants.OrganizationDescription)
	ImageStatusCmd.Flags().BoolVarP(&imageStatusWide, "wide", "w", false, "Show full digests")
	ImageStatusCmd.MarkFlagRequired(constants.OrganizationFlag)
}
