package model

// VM desired/observed models (Phase 4). JSON field names follow the protobuf
// JSON mapping used by the gateway (grpcurl): camelCase, repeated DesiredVm
// under "vms", repeated VmStatus under "statuses".

type VMDesired struct {
	Name     string `json:"name" yaml:"name"`
	Qcow2URL string `json:"qcow2Url" yaml:"qcow2Url"`
	Digest   string `json:"digest" yaml:"digest"`
	VCPUs    uint32 `json:"vcpus,omitempty" yaml:"vcpus,omitempty"`
	MemoryMB uint32 `json:"memoryMb,omitempty" yaml:"memoryMb,omitempty"`
}

type VMDesiredSetRequest struct {
	HostID string      `json:"hostId"`
	VMs    []VMDesired `json:"vms"`
}

type VMDesiredSetResponse struct {
	VMs []VMDesired `json:"vms"`
}

type VMDesiredGetRequest struct {
	HostID string `json:"hostId"`
}

type VMDesiredGetResponse struct {
	VMs []VMDesired `json:"vms"`
}

type VMStatus struct {
	Name     string `json:"name"`
	State    string `json:"state"`
	BootedOS string `json:"bootedOs,omitempty"`
}

type VMStatusGetRequest struct {
	HostID string `json:"hostId"`
}

type VMStatusGetResponse struct {
	Statuses []VMStatus `json:"statuses"`
}
