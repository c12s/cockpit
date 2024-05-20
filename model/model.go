package model

import (
	"time"
)

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

type Label struct {
	Key   string      `json:"key"`
	Value interface{} `json:"value"`
}

type LabelInput struct {
	Label  Label  `json:"label"`
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

type NodeResponse struct {
	Node Node `json:"node"`
}

type NodesResponse struct {
	Nodes []Node `json:"nodes"`
}

type ClaimNodesResponse struct {
	Nodes []Node `json:"node"`
}

type HTTPRequestConfig struct {
	URL         string
	Method      string
	Headers     map[string]string
	RequestBody interface{}
	Response    interface{}
	Token       string
	Timeout     time.Duration
}
