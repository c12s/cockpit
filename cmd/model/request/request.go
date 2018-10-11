package request

type MutateRequest struct {
	Request string   `json:"request"`
	Kind    string   `json:"kind"`
	MTData  Metadata `json:"metadata"`
	Regions []Region `json:"regions"`
}

type Metadata struct {
	Version      string `json:"version"`
	TaskName     string `json:"taskName"`
	Timestamp    int64  `json:"timestamp"`
	Namespace    string `json:"namespace"`
	ForceNSQueue bool   `json:"forceNamespaceQueue"`
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
	Request string            `json:"request"`
	Kind    string            `json:"kind"`
	MTData  Metadata          `json:"metadata"`
	Labels  map[string]string `json:"labels"`
}
