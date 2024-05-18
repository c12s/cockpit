package model

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type RegistrationDetails struct {
	Email    string `json:"email"`
	Name     string `json:"name"`
	Org      string `json:"org"`
	Password string `json:"password"`
	Surname  string `json:"surname"`
	Username string `json:"username"`
}

type TokenResponse struct {
	Token string `json:"token"`
}

type LabelInput struct {
	Label struct {
		Key   string      `json:"key"`
		Value interface{} `json:"value,omitempty"`
	} `json:"label"`
	NodeID string `json:"nodeId"`
	Org    string `json:"org"`
}

type DeleteLabelInput struct {
	LabelKey string `json:"labelKey"`
	NodeID   string `json:"nodeId"`
	Org      string `json:"org"`
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
	ID     string `json:"id"`
	Org    string `json:"org"`
	Labels []struct {
		Key   string `json:"key"`
		Value string `json:"value"`
	} `json:"labels"`
}

type NodesResponse struct {
	Nodes []Node `json:"nodes"`
}
