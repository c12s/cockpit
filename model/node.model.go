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
	ID     string  `json:"id"`
	Org    string  `json:"org"`
	Labels []Label `json:"labels"`
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
