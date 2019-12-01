package request

type MutateRequest struct {
	Version string   `json:"version"`
	Request string   `json:"request"`
	Kind    string   `json:"kind"`
	MTData  Metadata `json:"metadata"`
	Regions []Region `json:"regions"`
}

type Metadata struct {
	TaskName     string `json:"taskName"`
	Timestamp    int64  `json:"timestamp"`
	Namespace    string `json:"namespace"`
	ForceNSQueue bool   `json:"forceNamespaceQueue"`
	Queue        string `json:"queue"`
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
	Kind    string            `json:"kind"`
	Content map[string]string `json:"content"`
	Index   []string          `json:"index"`
}

type Strategy struct {
	Type     string            `json:"type"`
	Kind     string            `json:"kind"`
	Interval string            `json:"interval"`
	Retry    map[string]string `json:"retry"`
}

type Selector struct {
	Labels  map[string]string `json:"labels"`
	Compare map[string]string `json:"compare"`
}

type NMutateRequest struct {
	Version string            `json:"version"`
	Request string            `json:"request"`
	Kind    string            `json:"kind"`
	MTData  Metadata          `json:"metadata"`
	Name    string            `json:"name"`
	Labels  map[string]string `json:"labels"`
}

type NSResponse struct {
	Result []NSData `json:"data"`
}

type NSData struct {
	Age       string `json:"age"`
	Labels    string `json:"labels"`
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
}

type ConfigResponse struct {
	Result []ConfigData `json:"data"`
}

type ConfigData struct {
	RegionId  string `json:"regionId"`
	ClusterId string `json:"clusterId"`
	NodeId    string `json:"nodeId"`
	Configs   string `json:"configs"`
}

type ActionsResponse struct {
	Result []ActionsData `json:"data"`
}

type ActionsData struct {
	RegionId  string            `json:"regionId"`
	ClusterId string            `json:"clusterId"`
	NodeId    string            `json:"nodeId"`
	Actions   map[string]string `json:"actions"`
}

type SecretsResponse struct {
	Result []SecretsData `json:"data"`
}

type SecretsData struct {
	RegionId  string `json:"regionId"`
	ClusterId string `json:"clusterId"`
	NodeId    string `json:"nodeId"`
	Secrets   string `json:"secrets"`
}

type Credentials struct {
	Password string `json:"password"`
	Username string `json:"username"`
}

type SpanContext struct {
	TraceId       string            `json:"traceId"`
	SpanId        int64             `json:"spanId"`
	ParrentSpanId int64             `json:"spanId"`
	Baggage       map[string]string `json:"baggage"`
}

type Span struct {
	Context   SpanContext       `json:"spanContext"`
	Name      string            `json:"name"`
	Logs      map[string]string `json:"logs"`
	Tags      map[string]string `json:"tags"`
	StartTime int64             `json:"startTime"`
	EndTime   int64             `json:"endTime"`
}

type Trace struct {
	TraceId string `json:"traceId"`
	Trace   []Span `json:"trace"`
}

type Traces struct {
	Traces []Trace `json:"traces"`
}
