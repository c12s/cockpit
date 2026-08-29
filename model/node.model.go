package model

type ClaimNodesRequest struct {
	Org   string      `json:"org,omitempty"`
	Query []NodeQuery `json:"query,omitempty"`
}

type NodeQuery struct {
	LabelKey string      `json:"labelKey"`
	ShouldBe string      `json:"shouldBe"`
	Value    interface{} `json:"value"`
}

type Node struct {
	ID           string    `json:"id"`
	Org          string    `json:"org"`
	Labels       []Label   `json:"labels"`
	Resources    map[string]float64 `json:"resources,omitempty"`
	Capabilities []string  `json:"capabilities,omitempty"`
	OsStatus     *OsStatus `json:"osStatus,omitempty"`
}

type OsStatus struct {
	BootedImage   string     `json:"bootedImage,omitempty"`
	BootedDigest  string     `json:"bootedDigest,omitempty"`
	BootedVersion string     `json:"bootedVersion,omitempty"`
	Phase         string     `json:"phase,omitempty"`
	LastAttempt   *OsAttempt `json:"lastAttempt,omitempty"`
}

type OsAttempt struct {
	Image     string `json:"image,omitempty"`
	Digest    string `json:"digest,omitempty"`
	Result    string `json:"result,omitempty"`
	Timestamp int64  `json:"timestamp,omitempty"`
}

type NodeResponse struct {
	Node Node `json:"node"`
}

type NodesResponse struct {
	Nodes []Node `json:"nodes"`
}

type ClaimNodesResponse struct {
	Nodes []Node `json:"node"`
}

type OSImageUpdateRequest struct {
	NodeID       string `json:"nodeId"`
	ImageRef     string `json:"imageRef,omitempty"`
	Digest       string `json:"digest,omitempty"`
	RebootPolicy string `json:"rebootPolicy,omitempty"`
	Rollback     bool   `json:"rollback,omitempty"`
}

type OSImageUpdateResponse struct {
	Accepted bool `json:"accepted"`
}
