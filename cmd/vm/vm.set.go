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
	"gopkg.in/yaml.v3"
)

var (
	vmSetOrg  string
	vmSetHost string
	vmSetSpec string
	vmSetResp model.VMDesiredSetResponse
)

var VmsSetCmd = &cobra.Command{
	Use:     "set",
	Aliases: []string{"declare"},
	Short:   "Declare the desired VM set on a host node",
	Long:    "Store the desired VM set on the host node's record. The host's Star reconciler creates, restarts and destroys VMs to match it.",
	Run:     executeVmsSet,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return utils.ValidateRequiredFlags(cmd, []string{constants.HostFlag, constants.SpecFlag, constants.OrganizationFlag})
	},
}

func executeVmsSet(cmd *cobra.Command, args []string) {
	vms, err := parseVMSpec(vmSetSpec)
	if err != nil {
		fmt.Println("Error parsing VM spec:", err)
		os.Exit(1)
	}
	request := model.VMDesiredSetRequest{HostID: vmSetHost, VMs: vms}

	url := clients.BuildURL("core", "v1", "SetDesiredVms")
	token, err := utils.ReadTokenFromFile()
	if err != nil {
		fmt.Println("Error reading token:", err)
		os.Exit(1)
	}
	host, found := nodesByID(token, "ListOrgOwnedNodes", model.ClaimNodesRequest{Org: vmSetOrg})[vmSetHost]
	if err := checkVMHost(host, found, len(vms)); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	err = utils.SendHTTPRequest(model.HTTPRequestConfig{
		URL:         url,
		Method:      "POST",
		Token:       token,
		RequestBody: request,
		Response:    &vmSetResp,
		Timeout:     30 * time.Second,
	})
	if err != nil {
		fmt.Println("Error setting desired VMs:", err)
		os.Exit(1)
	}

	fmt.Printf("Desired VMs on host %s (%d):\n", vmSetHost, len(vmSetResp.VMs))
	for _, vm := range vmSetResp.VMs {
		fmt.Printf("  %-20s %-40s %s\n", vm.Name, vm.Qcow2URL, vm.Digest)
	}
}

// checkVMHost refuses a VM list for a node that cannot run VMs: only a host
// reporting run_vms (libvirt + /dev/kvm) runs the VM reconciler, anywhere
// else the list would be stored and silently ignored. Clearing the list is
// always allowed.
func checkVMHost(host model.Node, found bool, vmCount int) error {
	if vmCount == 0 {
		return nil
	}
	if !found {
		return fmt.Errorf("host %s is not a node of org %s (claim it first, or check --org)", vmSetHost, vmSetOrg)
	}
	for _, c := range host.Capabilities {
		if c == "run_vms" {
			return nil
		}
	}
	return fmt.Errorf("host %s cannot run VMs (no run_vms capability: needs libvirt and /dev/kvm)", host.ID)
}

// parseVMSpec accepts either a path to a YAML/JSON file or an inline JSON
// document (the plan's examples use both). JSON is valid YAML, so one parser
// handles both, and the document may be a bare list or {vms: [...]}.
func parseVMSpec(raw string) ([]model.VMDesired, error) {
	data := []byte(raw)
	if info, err := os.Stat(raw); err == nil && !info.IsDir() {
		if data, err = os.ReadFile(raw); err != nil {
			return nil, err
		}
	}
	// JSON is valid YAML, so one parser handles both formats. The document
	// may be {vms: [...]} (an empty list removes all desired VMs) or a bare
	// list.
	var probe map[string]any
	if err := yaml.Unmarshal(data, &probe); err == nil {
		if _, ok := probe["vms"]; ok {
			var wrapper model.VMDesiredSetRequest
			if err := yaml.Unmarshal(data, &wrapper); err != nil {
				return nil, err
			}
			return validateVMSpec(wrapper.VMs)
		}
	}
	var vms []model.VMDesired
	if err := yaml.Unmarshal(data, &vms); err != nil {
		return nil, fmt.Errorf("expected a list of VMs or {vms: [...]}: %v", err)
	}
	return validateVMSpec(vms)
}

func validateVMSpec(vms []model.VMDesired) ([]model.VMDesired, error) {
	for i, vm := range vms {
		if vm.Name == "" {
			return nil, fmt.Errorf("vm %d: name is required", i)
		}
		if vm.Qcow2URL == "" {
			return nil, fmt.Errorf("vm %s: qcow2Url is required", vm.Name)
		}
		if err := utils.ValidateDigest(vm.Digest); err != nil {
			return nil, fmt.Errorf("vm %s: %w", vm.Name, err)
		}
	}
	return vms, nil
}

func init() {
	VmsCmd.AddCommand(VmsSetCmd)
	VmsSetCmd.Flags().StringVar(&vmSetHost, constants.HostFlag, "", constants.HostDescription)
	VmsSetCmd.Flags().StringVar(&vmSetSpec, constants.SpecFlag, "", constants.SpecDescription)
	VmsSetCmd.Flags().StringVarP(&vmSetOrg, constants.OrganizationFlag, constants.OrganizationShorthandFlag, "", "Organization the host belongs to (required)")
	VmsSetCmd.MarkFlagRequired(constants.HostFlag)
	VmsSetCmd.MarkFlagRequired(constants.SpecFlag)
}
