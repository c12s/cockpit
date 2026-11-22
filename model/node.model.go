package model

import (
	"strconv"
	"strings"
)

// FlexInt64 tolerates the gateway's grpcurl/protojson responses, which
// marshal proto int64 fields as JSON strings ("1788115200"), while our own
// tooling emits plain numbers.
type FlexInt64 int64

func (f *FlexInt64) UnmarshalJSON(data []byte) error {
	s := strings.Trim(strings.TrimSpace(string(data)), "\"")
	if s == "" || s == "null" {
		*f = 0
		return nil
	}
	v, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return err
	}
	*f = FlexInt64(v)
	return nil
}

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
	ID           string             `json:"id"`
	Org          string             `json:"org"`
	Labels       []Label            `json:"labels"`
	Resources    map[string]float64 `json:"resources,omitempty"`
	Capabilities []string           `json:"capabilities,omitempty"`
	OsStatus     *OsStatus          `json:"osStatus,omitempty"`
	DesiredOs    *DesiredOsImage    `json:"desiredOsImage,omitempty"`
	// Liveness view filled by Magnetar from the last heartbeat:
	// online | stale | offline | unknown.
	LastSeen FlexInt64 `json:"lastSeen,omitempty"`
	Liveness string    `json:"liveness,omitempty"`
}

type DesiredOsImage struct {
	ImageRef     string `json:"imageRef,omitempty"`
	Digest       string `json:"digest,omitempty"`
	RebootPolicy string `json:"rebootPolicy,omitempty"`
}

type OsStatus struct {
	BootedImage   string     `json:"bootedImage,omitempty"`
	BootedDigest  string     `json:"bootedDigest,omitempty"`
	BootedVersion string     `json:"bootedVersion,omitempty"`
	Phase         string     `json:"phase,omitempty"`
	LastAttempt   *OsAttempt `json:"lastAttempt,omitempty"`
	// the other deployment on disk, target of `image rollback`
	RollbackImage  string `json:"rollbackImage,omitempty"`
	RollbackDigest string `json:"rollbackDigest,omitempty"`
}

type OsAttempt struct {
	Image     string    `json:"image,omitempty"`
	Digest    string    `json:"digest,omitempty"`
	Result    string    `json:"result,omitempty"`
	Timestamp FlexInt64 `json:"timestamp,omitempty"`
}

type NodeResponse struct {
	Node Node `json:"node"`
}

type DeleteNodeInput struct {
	NodeID string `json:"nodeId"`
}

type DeleteNodeResponse struct{}

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
