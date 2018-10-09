package request

type MutateRequest struct {
	Request   string   `json:"request"`
	Kind      string   `json:"kind"`
	Name      string   `json:"name"`
	Timestamp int64    `json:"timestamp"`
	Namespace string   `json:"namespace"`
	Regions   []Region `json:"regions"`
}

type Region struct {
	ID       string    `json:"regionID"`
	Clusters []Cluster `json:"cluster"`
}

type Cluster struct {
	ID       string    `json:"clusterID"`
	Payload  []Payload `json:"payload"`
	Strategy Strategy  `json:"strategy"`
	Selector Selector  `json:"selector"`
}

type Payload struct {
	Kind    string   `json:"kind"`
	Content []string `json:"content"`
}

type Strategy struct {
	Type string `json:"type"`
	Kind string `json:"kind"`
}

type Selector struct {
	Labels  map[string]string `json:"labels"`
	Compare map[string]string `json:"compare"`
}

type NMutateRequest struct {
	Request   string            `json:"request"`
	Kind      string            `json:"kind"`
	Timestamp int64             `json:"timestamp"`
	Name      string            `json:"name"`
	Labels    map[string]string `json:"labels"`
}
