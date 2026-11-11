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

// rolloutConfig is the parameterised rollout used by both `image rollout`
// and `image apply`. The command flags map onto it; the apply spec maps its
// target/rollout blocks onto it.
type rolloutConfig struct {
	org     string
	canary  string
	query   string
	nodes   string
	ref     string
	digest  string
	timeout time.Duration
	poll    time.Duration
	soak    time.Duration
	noWait  bool
	force   bool
}

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
	runRollout(rolloutConfig{
		org:     rolloutOrg,
		canary:  rolloutCanary,
		query:   rolloutQuery,
		nodes:   rolloutNodes,
		ref:     rolloutRef,
		digest:  rolloutDigest,
		timeout: rolloutTimeout,
		poll:    rolloutPoll,
		soak:    rolloutSoak,
		noWait:  rolloutNoWait,
		// The standalone rollout command always sends the update; it is not
		// an idempotent apply. The apply command controls this via --force.
		force: true,
	})
}

// runRollout drives the canary gate and then the targets. It is shared by
// the standalone rollout command and the declarative apply command.
func runRollout(cfg rolloutConfig) {
	targets, err := rolloutTargets(cfg)
	if err != nil {
		rolloutFail("resolving targets: %v", err)
	}
	if len(targets) == 0 {
		rolloutFail("no target nodes matched (the canary is never a target)")
	}
	rolloutLog("targets after canary: %s", strings.Join(targets, ", "))
	if err := rolloutGroup(cfg, cfg.canary, targets); err != nil {
		rolloutFail("%v", err)
	}
}

// rolloutGroup updates the canary, waits for its health gate and soaks it,
// then updates the targets and waits for them. Nodes whose desired image
// already matches are not re-sent (unless cfg.force). With no targets the
// soak is skipped: it only protects targets. Shared by `image rollout` and
// `image apply`.
func rolloutGroup(cfg rolloutConfig, canary string, targets []string) error {
	var nodes map[string]model.Node
	var err error
	if !cfg.force {
		nodes, err = orgNodesFor(cfg.org)
		if err != nil {
			return fmt.Errorf("fetching nodes: %w", err)
		}
	}
	matches := func(id string) bool {
		node, ok := nodes[id]
		return !cfg.force && ok && desiredMatches(node.DesiredOs, model.HostImageSource{Ref: cfg.ref, Digest: cfg.digest})
	}

	if matches(canary) {
		rolloutLog("canary %s already on %s; skipping update and checking health", canary, shortDigest(cfg.digest))
	} else {
		rolloutLog("canary %s: sending %s@%s", canary, cfg.ref, cfg.digest)
		if err := sendRolloutUpdate(cfg, canary); err != nil {
			return fmt.Errorf("canary %s: %w", canary, err)
		}
	}

	verdicts, err := waitForGate(cfg, []string{canary})
	if err != nil {
		return fmt.Errorf("canary %s: %w", canary, err)
	}
	if v := verdicts[canary]; v != gateHealthy {
		return fmt.Errorf("canary %s %s: rollout stopped, no target was updated", canary, v)
	}
	if len(targets) == 0 {
		rolloutLog("%s healthy on %s", canary, shortDigest(cfg.digest))
		return nil
	}
	rolloutLog("canary %s healthy on %s, soaking for %s", canary, shortDigest(cfg.digest), cfg.soak)
	if err := soak(cfg, canary); err != nil {
		return fmt.Errorf("canary %s: %w: rollout stopped, no target was updated", canary, err)
	}
	rolloutLog("canary %s stayed healthy: proceeding", canary)

	for _, id := range targets {
		if matches(id) {
			rolloutLog("target %s: already on %s, skipping", id, shortDigest(cfg.digest))
			continue
		}
		if err := sendRolloutUpdate(cfg, id); err != nil {
			return fmt.Errorf("target %s: %w (targets before it were already updated)", id, err)
		}
		rolloutLog("target %s: update accepted", id)
	}
	if cfg.noWait {
		rolloutLog("not waiting for targets; track with: cockpit image status --org %s", cfg.org)
		return nil
	}

	verdicts, err = waitForGate(cfg, targets)
	if err != nil {
		return fmt.Errorf("targets: %w", err)
	}
	failed := 0
	for _, id := range targets {
		rolloutLog("target %s: %s", id, verdicts[id])
		if verdicts[id] != gateHealthy {
			failed++
		}
	}
	if failed > 0 {
		return fmt.Errorf("%d of %d targets did not become healthy", failed, len(targets))
	}
	rolloutLog("rollout complete: canary + %d targets healthy on %s", len(targets), shortDigest(cfg.digest))
	return nil
}

// rolloutAll updates every node at once (no canary) and waits until all of
// them are healthy; any rollback or timeout fails the group.
func rolloutAll(cfg rolloutConfig, ids []string) error {
	for _, id := range ids {
		if err := sendRolloutUpdate(cfg, id); err != nil {
			return fmt.Errorf("node %s: %w (nodes before it were already updated)", id, err)
		}
		rolloutLog("%s: update accepted", id)
	}
	verdicts, err := waitForGate(cfg, ids)
	if err != nil {
		return err
	}
	failed := 0
	for _, id := range ids {
		rolloutLog("%s: %s", id, verdicts[id])
		if verdicts[id] != gateHealthy {
			failed++
		}
	}
	if failed > 0 {
		return fmt.Errorf("%d of %d nodes did not become healthy", failed, len(ids))
	}
	return nil
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
func waitForGate(cfg rolloutConfig, ids []string) (map[string]gateVerdict, error) {
	verdicts := make(map[string]gateVerdict, len(ids))
	last := make(map[string]string, len(ids))
	deadline := time.Now().Add(cfg.timeout)
	for {
		nodes, err := orgNodesFor(cfg.org)
		if err != nil {
			return nil, err
		}
		for _, id := range ids {
			if _, done := verdicts[id]; done {
				continue
			}
			node, ok := nodes[id]
			if !ok {
				return nil, fmt.Errorf("node %s is not in org %s", id, cfg.org)
			}
			if phase := phaseOf(node); phase != last[id] {
				rolloutLog("%s: phase %s", id, phase)
				last[id] = phase
			}
			if v := evaluateGate(node, cfg.digest); v != gatePending {
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
		time.Sleep(cfg.poll)
	}
}

// soak requires the node to stay healthy on the digest for cfg.soak: a
// single healthy report is not enough, a deployment can still be rejected
// (greenboot reboot, rollback) after it first looks fine.
func soak(cfg rolloutConfig, id string) error {
	deadline := time.Now().Add(cfg.soak)
	for time.Now().Before(deadline) {
		time.Sleep(cfg.poll)
		nodes, err := orgNodesFor(cfg.org)
		if err != nil {
			return err
		}
		node, ok := nodes[id]
		if !ok {
			return fmt.Errorf("node left org %s", cfg.org)
		}
		if v := evaluateGate(node, cfg.digest); v != gateHealthy {
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

func rolloutTargets(cfg rolloutConfig) ([]string, error) {
	var candidates []model.Node
	if cfg.nodes != "" {
		all, err := orgNodesFor(cfg.org)
		if err != nil {
			return nil, err
		}
		for _, id := range strings.Split(cfg.nodes, ",") {
			id = strings.TrimSpace(id)
			node, ok := all[id]
			if !ok {
				return nil, fmt.Errorf("node %s is not in org %s", id, cfg.org)
			}
			candidates = append(candidates, node)
		}
	} else {
		queries, err := utils.CreateNodeQuery(cfg.query)
		if err != nil {
			return nil, err
		}
		var resp model.NodesResponse
		if err := gatewayGet("QueryOrgOwnedNodes", model.ClaimNodesRequest{Org: cfg.org, Query: queries}, &resp); err != nil {
			return nil, err
		}
		candidates = resp.Nodes
	}
	var ids []string
	for _, node := range candidates {
		if node.ID == cfg.canary {
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

func orgNodesFor(org string) (map[string]model.Node, error) {
	var resp model.NodesResponse
	if err := gatewayGet("ListOrgOwnedNodes", model.ClaimNodesRequest{Org: org}, &resp); err != nil {
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
func sendRolloutUpdate(cfg rolloutConfig, nodeID string) error {
	for attempt := 0; ; attempt++ {
		accepted, err := requestImageUpdate(nodeID, cfg.ref, cfg.digest, "")
		if isRateLimited(err) && attempt < rateLimitRetries {
			time.Sleep(rateLimitBackoff)
			continue
		}
		if err != nil {
			if strings.Contains(err.Error(), "permission") {
				return fmt.Errorf("%w (grant os.update on the node's org, or on node %s)", err, nodeID)
			}
			return err
		}
		if !accepted {
			return errors.New("node did not acknowledge the update (offline or not running star)")
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
