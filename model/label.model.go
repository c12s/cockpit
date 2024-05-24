package model

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
