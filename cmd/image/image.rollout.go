package cmd

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/c12s/cockpit/clients"
	"github.com/c12s/cockpit/constants"
	"github.com/c12s/cockpit/model"
	"github.com/c12s/cockpit/utils"
	"github.com/spf13/cobra"
)

const (
	osUpdateCapability = "os_update"
	phaseHealthy       = "healthy"
	resultRolledBack   = "rolled_back"
	rateLimitRetries   = 8
	rateLimitBackoff   = 15 * time.Second
)

var (
	rolloutOrg     string
	rolloutCanary  string
	rolloutQuery   string
	rolloutNodes   string
	rolloutRef     string
	rolloutDigest  string
	rolloutTimeout time.Duration
	rolloutPoll    time.Duration
	rolloutSoak    time.Duration
	rolloutNoWait  bool
)

var ImageRolloutCmd = &cobra.Command{
	Use:   "rollout",
	Short: "Roll an OS image out behind a canary health gate",
	Long: `Update one canary node first and wait for its health gate. If the canary
boots the image and stays healthy for --soak, the image is sent to the target
nodes (--query label selector or --nodes list, within --org). If the canary
rolls back, leaves healthy during the soak or times out, the rollout stops and
no target is touched.

  cockpit image rollout --org c12s --canary <node-id> --query 'role = prod' \
      --image-ref 192.168.122.1:5001/os/edge-base:v3.9 --digest sha256:...`,
	Run: executeImageRollout,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if err := utils.ValidateRequiredFlags(cmd, []string{constants.OrganizationFlag, "canary", constants.ImageRefFlag, constants.DigestFlag}); err != nil {
			return err
		}
		if (rolloutQuery == "") == (rolloutNodes == "") {
			return errors.New("exactly one of --query or --nodes selects the targets")
		}
		return nil
	},
}

func executeImageRollout(cmd *cobra.Command, args []string) {
	targets, err := rolloutTargets()
	if err != nil {
		rolloutFail("resolving targets: %v", err)
	}
	if len(targets) == 0 {
		rolloutFail("no target nodes matched (the canary is never a target)")
	}
	rolloutLog("targets after canary: %s", strings.Join(targets, ", "))

	rolloutLog("canary %s: sending %s@%s", rolloutCanary, rolloutRef, rolloutDigest)
	if err := sendRolloutUpdate(rolloutCanary); err != nil {
		rolloutFail("canary: %v", err)
	}

	verdicts, err := waitForGate([]string{rolloutCanary})
	if err != nil {
		rolloutFail("canary: %v", err)
	}
	if v := verdicts[rolloutCanary]; v != gateHealthy {
		rolloutFail("canary %s %s: rollout stopped, no target was updated", rolloutCanary, v)
	}
	rolloutLog("canary %s healthy on %s, soaking for %s", rolloutCanary, rolloutDigest, rolloutSoak)
	if err := soak(rolloutCanary); err != nil {
		rolloutFail("canary %s: %v: rollout stopped, no target was updated", rolloutCanary, err)
	}
	rolloutLog("canary %s stayed healthy: proceeding", rolloutCanary)

	for _, id := range targets {
		if err := sendRolloutUpdate(id); err != nil {
			rolloutFail("target %s: %v (targets before it were already updated)", id, err)
		}
		rolloutLog("target %s: update accepted", id)
	}
	if rolloutNoWait {
		rolloutLog("not waiting for targets; track with: cockpit image status --org %s", rolloutOrg)
		return
	}

	verdicts, err = waitForGate(targets)
	if err != nil {
		rolloutFail("targets: %v", err)
	}
	failed := 0
	for _, id := range targets {
		rolloutLog("target %s: %s", id, verdicts[id])
		if verdicts[id] != gateHealthy {
			failed++
		}
	}
	if failed > 0 {
		rolloutFail("%d of %d targets did not become healthy", failed, len(targets))
	}
	rolloutLog("rollout complete: canary + %d targets healthy on %s", len(targets), rolloutDigest)
}

type gateVerdict string

const (
	gatePending    gateVerdict = "pending"
	gateHealthy    gateVerdict = "healthy"
	gateRolledBack gateVerdict = "rolled back"
	gateTimedOut   gateVerdict = "timed out"
)

// evaluateGate decides one node's health gate for the rolled-out digest.
// Healthy needs the digest booted and the node healthy (greenboot passed);
// a rollback is only counted for an attempt at this digest, so an older
// failure on the node does not stop the rollout. Anything else (switching,
// a failed pull being retried) is still pending.
func evaluateGate(node model.Node, digest string) gateVerdict {
	st := node.OsStatus
	if st == nil {
		return gatePending
	}
	if st.BootedDigest == digest && st.Phase == phaseHealthy {
		return gateHealthy
	}
	if st.LastAttempt != nil && st.LastAttempt.Digest == digest && st.LastAttempt.Result == resultRolledBack {
		return gateRolledBack
	}
	return gatePending
}

// waitForGate polls the org's nodes until every id has a final verdict or
// the timeout passes (remaining nodes are then reported as timed out).
func waitForGate(ids []string) (map[string]gateVerdict, error) {
	verdicts := make(map[string]gateVerdict, len(ids))
	last := make(map[string]string, len(ids))
	deadline := time.Now().Add(rolloutTimeout)
	for {
		nodes, err := orgNodes()
		if err != nil {
			return nil, err
		}
		for _, id := range ids {
			if _, done := verdicts[id]; done {
				continue
			}
			node, ok := nodes[id]
			if !ok {
				return nil, fmt.Errorf("node %s is not in org %s", id, rolloutOrg)
			}
			if phase := phaseOf(node); phase != last[id] {
				rolloutLog("%s: phase %s", id, phase)
				last[id] = phase
			}
			if v := evaluateGate(node, rolloutDigest); v != gatePending {
				verdicts[id] = v
			}
		}
		if len(verdicts) == len(ids) {
			return verdicts, nil
		}
		if time.Now().After(deadline) {
			for _, id := range ids {
				if _, done := verdicts[id]; !done {
					verdicts[id] = gateTimedOut
				}
			}
			return verdicts, nil
		}
		time.Sleep(rolloutPoll)
	}
}

// soak requires the node to stay healthy on the digest for rolloutSoak: a
// single healthy report is not enough, a deployment can still be rejected
// (greenboot reboot, rollback) after it first looks fine.
func soak(id string) error {
	deadline := time.Now().Add(rolloutSoak)
	for time.Now().Before(deadline) {
		time.Sleep(rolloutPoll)
		nodes, err := orgNodes()
		if err != nil {
			return err
		}
		node, ok := nodes[id]
		if !ok {
			return fmt.Errorf("node left org %s", rolloutOrg)
		}
		if v := evaluateGate(node, rolloutDigest); v != gateHealthy {
			return fmt.Errorf("left healthy during soak (phase %s)", phaseOf(node))
		}
	}
	return nil
}

func phaseOf(node model.Node) string {
	if node.OsStatus == nil || node.OsStatus.Phase == "" {
		return "unknown"
	}
	return node.OsStatus.Phase
}

func rolloutTargets() ([]string, error) {
	var candidates []model.Node
	if rolloutNodes != "" {
		all, err := orgNodes()
		if err != nil {
			return nil, err
		}
		for _, id := range strings.Split(rolloutNodes, ",") {
			id = strings.TrimSpace(id)
			node, ok := all[id]
			if !ok {
				return nil, fmt.Errorf("node %s is not in org %s", id, rolloutOrg)
			}
			candidates = append(candidates, node)
		}
	} else {
		queries, err := utils.CreateNodeQuery(rolloutQuery)
		if err != nil {
			return nil, err
		}
		var resp model.NodesResponse
		if err := gatewayGet("QueryOrgOwnedNodes", model.ClaimNodesRequest{Org: rolloutOrg, Query: queries}, &resp); err != nil {
			return nil, err
		}
		candidates = resp.Nodes
	}
	var ids []string
	for _, node := range candidates {
		if node.ID == rolloutCanary {
			continue
		}
		if !hasCapability(node, osUpdateCapability) {
			rolloutLog("skipping %s: no %s capability", node.ID, osUpdateCapability)
			continue
		}
		ids = append(ids, node.ID)
	}
	return ids, nil
}

func hasCapability(node model.Node, capability string) bool {
	for _, c := range node.Capabilities {
		if c == capability {
			return true
		}
	}
	return false
}

func orgNodes() (map[string]model.Node, error) {
	var resp model.NodesResponse
	if err := gatewayGet("ListOrgOwnedNodes", model.ClaimNodesRequest{Org: rolloutOrg}, &resp); err != nil {
		return nil, err
	}
	nodes := make(map[string]model.Node, len(resp.Nodes))
	for _, node := range resp.Nodes {
		nodes[node.ID] = node
	}
	return nodes, nil
}

// gatewayGet retries while the gateway rate-limits the caller: a rollout
// polls for minutes, so it regularly runs into the per-user limit.
func gatewayGet(method string, body, response interface{}) error {
	token, err := utils.ReadTokenFromFile()
	if err != nil {
		return fmt.Errorf("reading token: %w", err)
	}
	for attempt := 0; ; attempt++ {
		err = utils.SendHTTPRequest(model.HTTPRequestConfig{
			URL:         clients.BuildURL("core", "v1", method),
			Method:      "GET",
			Token:       token,
			Headers:     map[string]string{"Content-Type": "application/json"},
			RequestBody: body,
			Response:    response,
			Timeout:     30 * time.Second,
		})
		if !isRateLimited(err) || attempt >= rateLimitRetries {
			return err
		}
		time.Sleep(rateLimitBackoff)
	}
}

// sendRolloutUpdate retries while the gateway rate-limits the caller.
func sendRolloutUpdate(nodeID string) error {
	for attempt := 0; ; attempt++ {
		accepted, err := requestImageUpdate(nodeID, rolloutRef, rolloutDigest, "")
		if isRateLimited(err) && attempt < rateLimitRetries {
			time.Sleep(rateLimitBackoff)
			continue
		}
		if err != nil {
			return err
		}
		if !accepted {
			return errors.New("node did not acknowledge the update (offline, or this digest already rolled back on it)")
		}
		return nil
	}
}

func isRateLimited(err error) bool {
	return err != nil && strings.Contains(strings.ToLower(err.Error()), "rate limit")
}

func rolloutLog(format string, args ...interface{}) {
	fmt.Printf("[%s] %s\n", time.Now().Format("15:04:05"), fmt.Sprintf(format, args...))
}

func rolloutFail(format string, args ...interface{}) {
	rolloutLog("ERROR: "+format, args...)
	os.Exit(1)
}

func init() {
	ImageCmd.AddCommand(ImageRolloutCmd)
	f := ImageRolloutCmd.Flags()
	f.StringVarP(&rolloutOrg, constants.OrganizationFlag, constants.OrganizationShorthandFlag, "", constants.OrganizationDescription)
	f.StringVar(&rolloutCanary, "canary", "", "Node ID updated first; its health gate decides the rollout (required)")
	f.StringVarP(&rolloutQuery, constants.QueryFlag, constants.QueryFlagShorthandFlag, "", "Label query selecting target nodes within the org")
	f.StringVar(&rolloutNodes, "nodes", "", "Comma-separated target node IDs (instead of --query)")
	f.StringVar(&rolloutRef, constants.ImageRefFlag, "", constants.ImageRefDescription)
	f.StringVar(&rolloutDigest, constants.DigestFlag, "", constants.DigestDescription)
	f.DurationVar(&rolloutTimeout, "timeout", 15*time.Minute, "Maximum wait for the canary, then for the targets")
	f.DurationVar(&rolloutPoll, "poll", 10*time.Second, "Status poll interval")
	f.DurationVar(&rolloutSoak, "soak", 2*time.Minute, "How long the canary must stay healthy before targets are updated")
	f.BoolVar(&rolloutNoWait, "no-wait", false, "Return after the targets accepted the update")
	for _, name := range []string{constants.OrganizationFlag, "canary", constants.ImageRefFlag, constants.DigestFlag} {
		ImageRolloutCmd.MarkFlagRequired(name)
	}
}
