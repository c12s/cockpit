package model

// HostImageSpec is the declarative fleet file accepted by
// `cockpit image apply --spec`. It names the nodes it manages by id and says
// which OS image each runs and, for VM hosts, which VMs exist and which OS
// each VM runs. Nodes that are not in the file are never touched. It is a
// client-side document: the control plane stores only the resulting per-node
// desired image and per-host VM set.
type HostImageSpec struct {
	APIVersion string            `json:"apiVersion,omitempty" yaml:"apiVersion,omitempty"`
	Kind       string            `json:"kind,omitempty" yaml:"kind,omitempty"`
	Metadata   HostImageMetadata `json:"metadata,omitempty" yaml:"metadata,omitempty"`
	Spec       HostImageSpecBody `json:"spec,omitempty" yaml:"spec,omitempty"`
}

type HostImageMetadata struct {
	Org  string `json:"org,omitempty" yaml:"org,omitempty"`
	Name string `json:"name,omitempty" yaml:"name,omitempty"`
}

type HostImageSpecBody struct {
	// Images are the named OS images nodes and VMs refer to.
	Images  map[string]HostImageSource `json:"images,omitempty" yaml:"images,omitempty"`
	Nodes   []HostImageNode            `json:"nodes,omitempty" yaml:"nodes,omitempty"`
	Rollout *HostImageRollout          `json:"rollout,omitempty" yaml:"rollout,omitempty"`
}

// HostImageSource is an OS image. Without a digest the tag is resolved to
// the digest it points to when the spec is applied.
type HostImageSource struct {
	Ref          string `json:"ref,omitempty" yaml:"ref,omitempty"`
	Digest       string `json:"digest,omitempty" yaml:"digest,omitempty"`
	RebootPolicy string `json:"rebootPolicy,omitempty" yaml:"rebootPolicy,omitempty"`
}

// HostImageNode is one managed node. Image is optional (absent: OS not
// managed). VMs is nil when the key is absent (VMs not touched) and an empty
// list for `vms: []` (all VMs removed).
type HostImageNode struct {
	ID    string         `json:"id" yaml:"id"`
	Name  string         `json:"name,omitempty" yaml:"name,omitempty"`
	Image string         `json:"image,omitempty" yaml:"image,omitempty"`
	VMs   *[]HostImageVM `json:"vms,omitempty" yaml:"vms,omitempty"`
}

// HostImageVM is one VM on a host. Disk is used only to create the VM; Image
// is the OS it runs afterwards (optional). Its node id is derived from the
// host id and the name.
type HostImageVM struct {
	Name     string `json:"name" yaml:"name"`
	Disk     string `json:"disk" yaml:"disk"`
	Digest   string `json:"digest,omitempty" yaml:"digest,omitempty"`
	VCPUs    uint32 `json:"vcpus,omitempty" yaml:"vcpus,omitempty"`
	MemoryMB uint32 `json:"memoryMb,omitempty" yaml:"memoryMb,omitempty"`
	Image    string `json:"image,omitempty" yaml:"image,omitempty"`
}

// HostImageRollout tunes the per-group canary gate.
type HostImageRollout struct {
	Soak    string `json:"soak,omitempty" yaml:"soak,omitempty"`
	Timeout string `json:"timeout,omitempty" yaml:"timeout,omitempty"`
	Poll    string `json:"poll,omitempty" yaml:"poll,omitempty"`
}
