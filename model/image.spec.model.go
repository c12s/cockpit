package model

// HostImageSpec is the declarative, Terraform-like OS spec accepted by
// `cockpit image apply --spec`. It is a client-side document: the control
// plane stores only the resulting per-node desired image.
//
// The document has two accepted forms:
//
//   - the canonical form: {apiVersion, kind, metadata, spec}
//   - the flat form: {org, image|build, target, rollout}
//
// The canonical form's `spec` wins when both are present; flat fields are
// filled in only where the canonical form leaves a field empty.
//
// The image source is either `image` (a pinned digest) or `build` (a
// Containerfile context). When both are present, the build runs and its
// digest is used; the image block is the fallback for the ref/digest.
type HostImageSpec struct {
	APIVersion string            `json:"apiVersion,omitempty" yaml:"apiVersion,omitempty"`
	Kind       string            `json:"kind,omitempty" yaml:"kind,omitempty"`
	Metadata   HostImageMetadata `json:"metadata,omitempty" yaml:"metadata,omitempty"`
	Spec       HostImageSpecBody `json:"spec,omitempty" yaml:"spec,omitempty"`

	// Flat form. Ignored when the canonical form supplies the same field.
	Org     string            `json:"org,omitempty" yaml:"org,omitempty"`
	Image   *HostImageSource  `json:"image,omitempty" yaml:"image,omitempty"`
	Build   *HostImageBuild   `json:"build,omitempty" yaml:"build,omitempty"`
	Target  *HostImageTarget  `json:"target,omitempty" yaml:"target,omitempty"`
	Rollout *HostImageRollout `json:"rollout,omitempty" yaml:"rollout,omitempty"`
}

// HostImageMetadata identifies the spec and defaults the build's name/tag.
type HostImageMetadata struct {
	Org     string `json:"org,omitempty" yaml:"org,omitempty"`
	Name    string `json:"name,omitempty" yaml:"name,omitempty"`
	Version string `json:"version,omitempty" yaml:"version,omitempty"`
}

// HostImageSpecBody is the canonical `spec` block.
type HostImageSpecBody struct {
	Image   *HostImageSource  `json:"image,omitempty" yaml:"image,omitempty"`
	Build   *HostImageBuild   `json:"build,omitempty" yaml:"build,omitempty"`
	Target  *HostImageTarget  `json:"target,omitempty" yaml:"target,omitempty"`
	Rollout *HostImageRollout `json:"rollout,omitempty" yaml:"rollout,omitempty"`
}

// HostImageSource is a pinned OS image. The digest is the immutable unit.
type HostImageSource struct {
	Ref          string `json:"ref,omitempty" yaml:"ref,omitempty"`
	Digest       string `json:"digest,omitempty" yaml:"digest,omitempty"`
	RebootPolicy string `json:"rebootPolicy,omitempty" yaml:"rebootPolicy,omitempty"`
}

// HostImageBuild describes a build context. The build runs first; its digest
// becomes the hand-off to the update path.
type HostImageBuild struct {
	Context       string            `json:"context,omitempty" yaml:"context,omitempty"`
	Containerfile string            `json:"containerfile,omitempty" yaml:"containerfile,omitempty"`
	Registry      string            `json:"registry,omitempty" yaml:"registry,omitempty"`
	Name          string            `json:"name,omitempty" yaml:"name,omitempty"`
	Tag           string            `json:"tag,omitempty" yaml:"tag,omitempty"`
	BuildArgs     map[string]string `json:"buildArgs,omitempty" yaml:"buildArgs,omitempty"`
	Timeout       string            `json:"timeout,omitempty" yaml:"timeout,omitempty"`
}

// HostImageTarget selects the nodes. Exactly one of Nodes or Query is set.
type HostImageTarget struct {
	Org   string   `json:"org,omitempty" yaml:"org,omitempty"`
	Nodes []string `json:"nodes,omitempty" yaml:"nodes,omitempty"`
	Query string   `json:"query,omitempty" yaml:"query,omitempty"`
}

// HostImageRollout is the canary/soak policy. When absent, targets are
// updated directly.
type HostImageRollout struct {
	Canary  string `json:"canary,omitempty" yaml:"canary,omitempty"`
	Soak    string `json:"soak,omitempty" yaml:"soak,omitempty"`
	Timeout string `json:"timeout,omitempty" yaml:"timeout,omitempty"`
	Poll    string `json:"poll,omitempty" yaml:"poll,omitempty"`
	NoWait  bool   `json:"noWait,omitempty" yaml:"noWait,omitempty"`
}

// Normalize merges the flat form into the canonical form. The canonical form
// wins when both carry a value.
func (s HostImageSpec) Normalize() HostImageSpecBody {
	body := s.Spec
	if body.Image == nil {
		body.Image = s.Image
	}
	if body.Build == nil {
		body.Build = s.Build
	}
	if body.Target == nil {
		body.Target = s.Target
	}
	if body.Rollout == nil {
		body.Rollout = s.Rollout
	}
	return body
}
