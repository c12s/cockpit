package cmd

import (
	"encoding/base64"
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
	"gopkg.in/yaml.v3"
)

var (
	applySpec    string
	applyDryRun  bool
	applyForce   bool
	applyTimeout time.Duration
)

var ImageApplyCmd = &cobra.Command{
	Use:     "apply",
	Aliases: []string{"declare"},
	Short:   "Apply a declarative HostImageSpec (build or pinned image, targets, rollout)",
	Long: `Apply a declarative host image spec. The spec is a client-side document:
the control plane only stores the resulting per-node desired image.

  spec.image is a pinned digest; spec.build is a Containerfile context.
  When both are present the build runs and its digest is used; spec.image
  is the fallback. The build goes through the gateway and Supernova, so the
  command waits for it to finish and returns the digest.

  spec.target selects the nodes by nodes or query (exactly one).
  spec.rollout runs a canary first, soaks, then updates the targets. Without
  spec.rollout, the targets are updated directly.

  --dry-run prints the plan without building or updating; --force re-applies
  even when the desired image already matches.`,
	Run: executeImageApply,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return utils.ValidateRequiredFlags(cmd, []string{constants.SpecFlag})
	},
}

func executeImageApply(cmd *cobra.Command, args []string) {
	spec, err := loadHostImageSpec(applySpec)
	if err != nil {
		applyFail("parsing spec: %v", err)
	}

	meta := spec.Metadata
	if meta.Org == "" {
		meta.Org = spec.Org
	}
	body := spec.Normalize()

	target, err := normalizeTarget(body.Target, meta)
	if err != nil {
		applyFail("validating target: %v", err)
	}
	rollout, err := normalizeRollout(body.Rollout)
	if err != nil {
		applyFail("validating rollout: %v", err)
	}
	source, err := normalizeImageSource(body.Image, body.Build)
	if err != nil {
		applyFail("validating image source: %v", err)
	}

	var build *model.HostImageBuild
	if body.Build != nil {
		build, err = normalizeBuild(body.Build, meta, source, applyTimeout)
		if err != nil {
			applyFail("validating build: %v", err)
		}
	}

	if applyDryRun {
		printApplyPlan(meta, target, rollout, source, build)
		return
	}

	if build != nil {
		source, err = buildHostImage(build, source, applyTimeout)
		if err != nil {
			applyFail("building image: %v", err)
		}
	}
	if err := utils.ValidateDigest(source.Digest); err != nil {
		applyFail("invalid image digest: %v", err)
	}

	if rollout != nil {
		soak, _ := time.ParseDuration(rollout.Soak)
		timeout, _ := time.ParseDuration(rollout.Timeout)
		poll, _ := time.ParseDuration(rollout.Poll)
		runRollout(rolloutConfig{
			org:     target.Org,
			canary:  rollout.Canary,
			query:   target.Query,
			nodes:   strings.Join(target.Nodes, ","),
			ref:     source.Ref,
			digest:  source.Digest,
			timeout: timeout,
			poll:    poll,
			soak:    soak,
			noWait:  rollout.NoWait,
			force:   applyForce,
		})
		return
	}

	if err := applyToTargets(source, target, applyForce); err != nil {
		applyFail("applying image: %v", err)
	}
}

// loadHostImageSpec reads a path or an inline JSON document. JSON is valid
// YAML, so one parser handles both.
func loadHostImageSpec(raw string) (model.HostImageSpec, error) {
	var spec model.HostImageSpec
	data := []byte(raw)
	if info, err := os.Stat(raw); err == nil && !info.IsDir() {
		var readErr error
		if data, readErr = os.ReadFile(raw); readErr != nil {
			return spec, readErr
		}
	}
	if err := yaml.Unmarshal(data, &spec); err != nil {
		return spec, fmt.Errorf("unmarshal YAML/JSON: %w", err)
	}
	return spec, nil
}

// normalizeImageSource validates the pinned-image or build source. A build
// may be the only source; its image block is optional fallback.
func normalizeImageSource(image *model.HostImageSource, build *model.HostImageBuild) (model.HostImageSource, error) {
	if image == nil && build == nil {
		return model.HostImageSource{}, errors.New("spec.image or spec.build is required")
	}
	src := model.HostImageSource{}
	if image != nil {
		src.Ref = image.Ref
		src.Digest = image.Digest
		src.RebootPolicy = image.RebootPolicy
		if src.RebootPolicy == "" {
			src.RebootPolicy = "immediate"
		}
		if src.Digest != "" {
			if err := utils.ValidateDigest(src.Digest); err != nil {
				return model.HostImageSource{}, err
			}
		}
	}
	if build == nil {
		if src.Ref == "" {
			return model.HostImageSource{}, errors.New("spec.image.ref is required")
		}
		if src.Digest == "" {
			return model.HostImageSource{}, errors.New("spec.image.digest is required")
		}
	}
	return src, nil
}

// normalizeBuild applies the defaults and validates the build block. The
// metadata name/version and the image ref fill in missing name/tag/registry.
func normalizeBuild(build *model.HostImageBuild, meta model.HostImageMetadata, image model.HostImageSource, defaultTimeout time.Duration) (*model.HostImageBuild, error) {
	if build == nil {
		return nil, nil
	}
	out := *build
	if out.Context == "" {
		out.Context = "."
	}
	if out.Containerfile == "" {
		out.Containerfile = "Containerfile"
	}
	if out.Registry == "" {
		out.Registry = registryFromRef(image.Ref)
		if out.Registry == "" {
			out.Registry = "192.168.122.1:5001"
		}
	}
	if out.Name == "" {
		out.Name = repoFromRef(image.Ref)
		if out.Name == "" {
			out.Name = meta.Name
		}
	}
	if out.Name == "" {
		return nil, errors.New("spec.build.name is required (or metadata.name or spec.image.ref)")
	}
	if out.Tag == "" {
		out.Tag = tagFromRef(image.Ref)
		if out.Tag == "" {
			out.Tag = meta.Version
		}
		if out.Tag == "" {
			out.Tag = "latest"
		}
	}
	if out.Timeout == "" {
		out.Timeout = defaultTimeout.String()
	}
	if _, err := time.ParseDuration(out.Timeout); err != nil {
		return nil, fmt.Errorf("spec.build.timeout: %w", err)
	}
	if out.BuildArgs == nil {
		out.BuildArgs = map[string]string{}
	}
	return &out, nil
}

// normalizeTarget validates the target and defaults the org to metadata.org.
func normalizeTarget(target *model.HostImageTarget, meta model.HostImageMetadata) (*model.HostImageTarget, error) {
	if target == nil {
		return nil, errors.New("spec.target is required")
	}
	org := target.Org
	if org == "" {
		org = meta.Org
	}
	if org == "" {
		return nil, errors.New("spec.target.org is required (or metadata.org)")
	}
	hasNodes := len(target.Nodes) > 0
	hasQuery := target.Query != ""
	if hasNodes == hasQuery {
		return nil, errors.New("spec.target.nodes and spec.target.query are mutually exclusive; set exactly one")
	}
	if hasNodes {
		seen := make(map[string]bool, len(target.Nodes))
		nodes := make([]string, 0, len(target.Nodes))
		for _, id := range target.Nodes {
			id = strings.TrimSpace(id)
			if id == "" || seen[id] {
				continue
			}
			seen[id] = true
			nodes = append(nodes, id)
		}
		if len(nodes) == 0 {
			return nil, errors.New("spec.target.nodes must not be empty")
		}
		target.Nodes = nodes
	}
	target.Org = org
	return target, nil
}

// normalizeRollout validates the rollout block and fills in its defaults.
func normalizeRollout(rollout *model.HostImageRollout) (*model.HostImageRollout, error) {
	if rollout == nil {
		return nil, nil
	}
	out := *rollout
	if out.Canary == "" {
		return nil, errors.New("spec.rollout.canary is required")
	}
	if out.Soak == "" {
		out.Soak = "2m"
	}
	if out.Timeout == "" {
		out.Timeout = "15m"
	}
	if out.Poll == "" {
		out.Poll = "10s"
	}
	for _, d := range []struct {
		name  string
		value string
	}{
		{"soak", out.Soak},
		{"timeout", out.Timeout},
		{"poll", out.Poll},
	} {
		if _, err := time.ParseDuration(d.value); err != nil {
			return nil, fmt.Errorf("spec.rollout.%s: %w", d.name, err)
		}
	}
	return &out, nil
}

// buildHostImage runs the build and returns the resulting source. The build's
// ref/digest win; the image block fills in anything the build omits.
func buildHostImage(build *model.HostImageBuild, image model.HostImageSource, defaultTimeout time.Duration) (model.HostImageSource, error) {
	timeout := defaultTimeout
	if build.Timeout != "" {
		d, err := time.ParseDuration(build.Timeout)
		if err != nil {
			return model.HostImageSource{}, fmt.Errorf("spec.build.timeout: %w", err)
		}
		timeout = d
	}
	tgz, err := utils.TarGzDir(build.Context)
	if err != nil {
		return model.HostImageSource{}, fmt.Errorf("packing context %s: %w", build.Context, err)
	}
	req := model.BuildImageRequest{
		Spec: model.BuildSpec{
			Name:          build.Name,
			Tag:           build.Tag,
			Containerfile: build.Containerfile,
			Registry:      build.Registry,
			BuildArgs:     build.BuildArgs,
		},
		Context: base64.StdEncoding.EncodeToString(tgz),
	}
	var job model.BuildJob
	if err := sendBuildPost(clients.BuildURL("core", "v1", "BuildImage"), req, &job); err != nil {
		return model.HostImageSource{}, fmt.Errorf("starting build: %w", err)
	}
	fmt.Printf("Build job %s started (%s:%s, %d bytes context).\n", job.JobID, build.Name, build.Tag, len(tgz))
	status, err := waitForBuild(job.JobID, timeout)
	if err != nil {
		printBuildJob(status)
		return model.HostImageSource{}, fmt.Errorf("waiting for build %s: %w", job.JobID, err)
	}
	printBuildJob(status)
	if status.Status == "failed" {
		return model.HostImageSource{}, fmt.Errorf("build %s failed: %s", job.JobID, status.Error)
	}
	ref := status.ImageRef
	if ref == "" {
		ref = image.Ref
	}
	digest := status.Digest
	if digest == "" {
		digest = image.Digest
	}
	if ref == "" {
		return model.HostImageSource{}, errors.New("build did not return an image ref")
	}
	if digest == "" {
		return model.HostImageSource{}, errors.New("build did not return a digest")
	}
	if image.Ref != "" && image.Ref != ref {
		fmt.Printf("Note: build image ref %s differs from spec.image.ref %s; using the build's ref.\n", ref, image.Ref)
	}
	if image.Digest != "" && image.Digest != digest {
		fmt.Printf("Note: build digest %s differs from spec.image.digest %s; using the build's digest.\n", shortDigest(digest), shortDigest(image.Digest))
	}
	policy := image.RebootPolicy
	if policy == "" {
		policy = "immediate"
	}
	return model.HostImageSource{Ref: ref, Digest: digest, RebootPolicy: policy}, nil
}

// applyToTargets updates every target directly, skipping nodes whose desired
// image already matches unless force is set.
func applyToTargets(source model.HostImageSource, target *model.HostImageTarget, force bool) error {
	targets, err := resolveTargetNodes(target)
	if err != nil {
		return err
	}
	if len(targets) == 0 {
		return errors.New("no target nodes matched")
	}
	nodes, err := orgNodesFor(target.Org)
	if err != nil {
		return err
	}
	var updated, skipped int
	for _, id := range targets {
		node, ok := nodes[id]
		if !ok {
			return fmt.Errorf("node %s is not in org %s", id, target.Org)
		}
		if !force && desiredMatches(node.DesiredOs, source) {
			fmt.Printf("Node %s already on %s; skipping.\n", id, shortDigest(source.Digest))
			skipped++
			continue
		}
		accepted, err := applyImageUpdate(id, source.Ref, source.Digest, source.RebootPolicy)
		if err != nil {
			return fmt.Errorf("node %s: %w", id, err)
		}
		if !accepted {
			fmt.Printf("Node %s did not acknowledge the update (offline?). The desired image is recorded; the node applies it when it reconnects.\n", id)
			updated++
			continue
		}
		fmt.Printf("Node %s accepted the update to %s@%s.\n", id, source.Ref, shortDigest(source.Digest))
		updated++
	}
	fmt.Printf("Applied %s@%s: %d updated, %d skipped.\n", source.Ref, shortDigest(source.Digest), updated, skipped)
	return nil
}

// resolveTargetNodes resolves the target selector to concrete node ids.
func resolveTargetNodes(target *model.HostImageTarget) ([]string, error) {
	if target == nil {
		return nil, errors.New("target is required")
	}
	hasNodes := len(target.Nodes) > 0
	hasQuery := target.Query != ""
	if hasNodes == hasQuery {
		return nil, errors.New("target.nodes and target.query are mutually exclusive; set exactly one")
	}
	if hasNodes {
		nodes, err := orgNodesFor(target.Org)
		if err != nil {
			return nil, err
		}
		seen := make(map[string]bool, len(target.Nodes))
		var ids []string
		for _, id := range target.Nodes {
			id = strings.TrimSpace(id)
			if id == "" || seen[id] {
				continue
			}
			node, ok := nodes[id]
			if !ok {
				return nil, fmt.Errorf("node %s is not in org %s", id, target.Org)
			}
			if !hasCapability(node, osUpdateCapability) {
				fmt.Printf("Skipping node %s: no %s capability.\n", id, osUpdateCapability)
				continue
			}
			seen[id] = true
			ids = append(ids, id)
		}
		return ids, nil
	}
	queries, err := utils.CreateNodeQuery(target.Query)
	if err != nil {
		return nil, err
	}
	var resp model.NodesResponse
	if err := gatewayGet("QueryOrgOwnedNodes", model.ClaimNodesRequest{Org: target.Org, Query: queries}, &resp); err != nil {
		return nil, err
	}
	seen := make(map[string]bool, len(resp.Nodes))
	var ids []string
	for _, node := range resp.Nodes {
		if seen[node.ID] {
			continue
		}
		if !hasCapability(node, osUpdateCapability) {
			fmt.Printf("Skipping node %s: no %s capability.\n", node.ID, osUpdateCapability)
			continue
		}
		seen[node.ID] = true
		ids = append(ids, node.ID)
	}
	return ids, nil
}

// desiredMatches reports whether the node's declared desired image already
// points at the source. The digest is the immutable unit; a matching ref or
// policy is only required when the source carries one.
func desiredMatches(desired *model.DesiredOsImage, source model.HostImageSource) bool {
	if desired == nil || desired.Digest != source.Digest {
		return false
	}
	if source.Ref != "" && desired.ImageRef != "" && desired.ImageRef != source.Ref {
		return false
	}
	if source.RebootPolicy != "" && desired.RebootPolicy != "" && desired.RebootPolicy != source.RebootPolicy {
		return false
	}
	return true
}

// applyImageUpdate retries while the gateway rate-limits the caller.
func applyImageUpdate(nodeID, ref, digest, policy string) (bool, error) {
	for attempt := 0; ; attempt++ {
		accepted, err := requestImageUpdate(nodeID, ref, digest, policy)
		if isRateLimited(err) && attempt < rateLimitRetries {
			time.Sleep(rateLimitBackoff)
			continue
		}
		return accepted, err
	}
}

// registryFromRef extracts the registry from an image reference when the
// first component looks like a host (a dot, a port, or localhost).
func registryFromRef(ref string) string {
	i := strings.Index(ref, "/")
	if i <= 0 {
		return ""
	}
	host := ref[:i]
	if strings.Contains(host, ".") || strings.Contains(host, ":") || host == "localhost" {
		return host
	}
	return ""
}

// repoFromRef extracts the repository without the registry and tag.
func repoFromRef(ref string) string {
	host := registryFromRef(ref)
	repo := ref
	if host != "" {
		repo = strings.TrimPrefix(ref, host+"/")
	}
	if i := strings.LastIndex(repo, ":"); i > strings.LastIndex(repo, "/") {
		repo = repo[:i]
	}
	return repo
}

// tagFromRef extracts the tag when the last colon is after the last slash.
func tagFromRef(ref string) string {
	i := strings.LastIndex(ref, ":")
	if i > strings.LastIndex(ref, "/") {
		return ref[i+1:]
	}
	return ""
}

func printApplyPlan(meta model.HostImageMetadata, target *model.HostImageTarget, rollout *model.HostImageRollout, source model.HostImageSource, build *model.HostImageBuild) {
	fmt.Printf("HostImageSpec: org=%s name=%s version=%s\n", meta.Org, meta.Name, meta.Version)
	if build != nil {
		fmt.Printf("Build: %s:%s registry=%s context=%s containerfile=%s timeout=%s\n", build.Name, build.Tag, build.Registry, build.Context, build.Containerfile, build.Timeout)
		if source.Ref != "" || source.Digest != "" {
			fmt.Printf("Image fallback: %s@%s\n", source.Ref, source.Digest)
		}
	} else {
		fmt.Printf("Image: %s@%s rebootPolicy=%s\n", source.Ref, source.Digest, source.RebootPolicy)
	}
	if len(target.Nodes) > 0 {
		fmt.Printf("Target: org=%s nodes=%s\n", target.Org, strings.Join(target.Nodes, ","))
	} else {
		fmt.Printf("Target: org=%s query=%q\n", target.Org, target.Query)
	}
	if rollout != nil {
		fmt.Printf("Rollout: canary=%s soak=%s timeout=%s poll=%s noWait=%t\n", rollout.Canary, rollout.Soak, rollout.Timeout, rollout.Poll, rollout.NoWait)
	} else {
		fmt.Println("Rollout: none (apply targets directly)")
	}
	fmt.Println("Dry run: no build or update was performed.")
}

func applyFail(format string, args ...interface{}) {
	fmt.Printf("Error: "+format+"\n", args...)
	os.Exit(1)
}

func init() {
	ImageCmd.AddCommand(ImageApplyCmd)
	f := ImageApplyCmd.Flags()
	f.StringVar(&applySpec, constants.SpecFlag, "", constants.HostImageSpecDescription)
	f.BoolVar(&applyDryRun, "dry-run", false, "Print the plan and exit without building or updating")
	f.BoolVar(&applyForce, "force", false, "Update nodes even when the desired image already matches the digest")
	f.DurationVar(&applyTimeout, "timeout", 30*time.Minute, "Deadline for a build started by the spec")
	ImageApplyCmd.MarkFlagRequired(constants.SpecFlag)
}
