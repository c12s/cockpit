package model

type Relation struct {
	From Entity `json:"from"`
	To   Entity `json:"to"`
}

type Entity struct {
	ID   string `json:"id"`
	Kind string `json:"kind"`
}
