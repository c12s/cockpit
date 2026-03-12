package model

type GetStarmapChartByMetadataReq struct {
	Maintainer    string `json:"maintainer" yaml:"maintainer"`
	Name          string `json:"name" yaml:"name"`
	Namespace     string `json:"namespace,omitempty" yaml:"namespace,omitempty"`
	SchemaVersion string `json:"schemaVersion,omitempty" yaml:"schemaVersion,omitempty"`
}

type GetStarmapChartByIdReq struct {
	ChartId       string `json:"chartId" yaml:"chartId"`
	Maintainer    string `json:"maintainer" yaml:"maintainer"`
	Namespace     string `json:"namespace,omitempty" yaml:"namespace,omitempty"`
	SchemaVersion string `json:"schemaVersion,omitempty" yaml:"schemaVersion,omitempty"`
}

type GetStarmapChartsByLabelsReq struct {
	Maintainer string            `json:"maintainer" yaml:"maintainer"`
	Namespace  string            `json:"namespace,omitempty" yaml:"namespace,omitempty"`
	Labels     map[string]string `json:"labels,omitempty" yaml:"labels,omitempty"`
}

type GetStarmapChartMissingLayersReq struct {
	ChartId       string   `json:"chartId" yaml:"chartId"`
	Maintainer    string   `json:"maintainer" yaml:"maintainer"`
	Namespace     string   `json:"namespace,omitempty" yaml:"namespace,omitempty"`
	Layers        []string `json:"layers,omitempty" yaml:"layers,omitempty"`
	SchemaVersion string   `json:"schemaVersion,omitempty" yaml:"schemaVersion,omitempty"`
}

type GetMissingLayersResp struct {
	ChartId          string                      `json:"chartId" yaml:"chartId"`
	Maintainer       string                      `json:"maintainer" yaml:"maintainer"`
	Namespace        string                      `json:"namespace" yaml:"namespace"`
	ApiVersion       string                      `json:"apiVersion" yaml:"apiVersion"`
	SchemaVersion    string                      `json:"schemaVersion" yaml:"schemaVersion"`
	DataSources      map[string]*DataSource      `json:"dataSources,omitempty" yaml:"dataSources,omitempty"`
	StoredProcedures map[string]*StoredProcedure `json:"storedProcedures,omitempty" yaml:"storedProcedures,omitempty"`
	EventTriggers    map[string]*EventTrigger    `json:"eventTriggers,omitempty" yaml:"eventTriggers,omitempty"`
	Events           map[string]*Event           `json:"events,omitempty" yaml:"events,omitempty"`
}

type DeleteChartReq struct {
	Id            string `json:"id" yaml:"id"`
	Maintainer    string `json:"maintainer" yaml:"maintainer"`
	Namespace     string `json:"namespace" yaml:"namespace"`
	Name          string `json:"name" yaml:"name"`
	Kind          string `json:"kind" yaml:"kind"`
	SchemaVersion string `json:"schemaVersion" yaml:"schemaVersion"`
}

type GetChartsLabelsResp struct {
	Charts []GetChartResp
}

type SwitchCheckpointReq struct {
	ChartId    string   `json:"chartId" yaml:"chartId"`
	Maintainer string   `json:"maintainer" yaml:"maintainer"`
	Namespace  string   `json:"namespace" yaml:"namespace"`
	OldVersion string   `json:"oldVersion" yaml:"oldVersion"`
	NewVersion string   `json:"newVersion" yaml:"newVersion"`
	Layers     []string `json:"layers,omitempty" yaml:"layers,omitempty"`
}

type LayersResp struct {
	DataSources      map[string]*DataSource      `json:"dataSources,omitempty" yaml:"dataSources,omitempty"`
	StoredProcedures map[string]*StoredProcedure `json:"storedProcedures,omitempty" yaml:"storedProcedures,omitempty"`
	EventTriggers    map[string]*EventTrigger    `json:"eventTriggers,omitempty" yaml:"eventTriggers,omitempty"`
	Events           map[string]*Event           `json:"events,omitempty" yaml:"events,omitempty"`
}

type SwitchCheckpointResp struct {
	Start struct {
		DataSources      map[string]*DataSource      `json:"dataSources,omitempty" yaml:"dataSources,omitempty"`
		StoredProcedures map[string]*StoredProcedure `json:"storedProcedures,omitempty" yaml:"storedProcedures,omitempty"`
		EventTriggers    map[string]*EventTrigger    `json:"eventTriggers,omitempty" yaml:"eventTriggers,omitempty"`
		Events           map[string]*Event           `json:"events,omitempty" yaml:"events,omitempty"`
	} `json:"start" yaml:"start"`

	Stop struct {
		DataSources      map[string]*DataSource      `json:"dataSources,omitempty" yaml:"dataSources,omitempty"`
		StoredProcedures map[string]*StoredProcedure `json:"storedProcedures,omitempty" yaml:"storedProcedures,omitempty"`
		EventTriggers    map[string]*EventTrigger    `json:"eventTriggers,omitempty" yaml:"eventTriggers,omitempty"`
		Events           map[string]*Event           `json:"events,omitempty" yaml:"events,omitempty"`
	} `json:"stop" yaml:"stop"`

	Download struct {
		DataSources      map[string]*DataSource      `json:"dataSources,omitempty" yaml:"dataSources,omitempty"`
		StoredProcedures map[string]*StoredProcedure `json:"storedProcedures,omitempty" yaml:"storedProcedures,omitempty"`
		EventTriggers    map[string]*EventTrigger    `json:"eventTriggers,omitempty" yaml:"eventTriggers,omitempty"`
		Events           map[string]*Event           `json:"events,omitempty" yaml:"events,omitempty"`
	} `json:"download" yaml:"download"`
}

type GetStarmapChartATimelineReq struct {
	ChartId    string `json:"chartId" yaml:"chartId"`
	Maintainer string `json:"maintainer" yaml:"maintainer"`
	Namespace  string `json:"namespace" yaml:"namespace"`
}

type PutChartResp struct {
	ApiVersion    string `json:"apiVersion" yaml:"apiVersion"`
	SchemaVersion string `json:"schemaVersion" yaml:"schemaVersion"`
	Id            string `json:"id" yaml:"id"`
	Name          string `json:"name" yaml:"name"`
	Namespace     string `json:"namespace" yaml:"namespace"`
	Maintainer    string `json:"maintainer" yaml:"maintainer"`
	Kind          string `json:"kind" yaml:"kind"`
}

type GetChartResp struct {
	ApiVersion    string `json:"apiVersion" yaml:"apiVersion"`
	SchemaVersion string `json:"schemaVersion" yaml:"schemaVersion"`

	Metadata struct {
		Id          string            `json:"id" yaml:"id"`
		Name        string            `json:"name" yaml:"name"`
		Namespace   string            `json:"namespace" yaml:"namespace"`
		Maintainer  string            `json:"maintainer" yaml:"maintainer"`
		Description string            `json:"description" yaml:"description"`
		Visibility  string            `json:"visibility" yaml:"visibility"`
		Engine      string            `json:"engine" yaml:"engine"`
		Labels      map[string]string `json:"labels,omitempty" yaml:"labels,omitempty"`
	} `json:"metadata" yaml:"metadata"`

	Chart Chart `json:"chart" yaml:"chart"`
}

type Chart struct {
	DataSources      map[string]*DataSource      `json:"dataSources,omitempty" yaml:"dataSources,omitempty"`
	StoredProcedures map[string]*StoredProcedure `json:"storedProcedures,omitempty" yaml:"storedProcedures,omitempty"`
	EventTriggers    map[string]*EventTrigger    `json:"eventTriggers,omitempty" yaml:"eventTriggers,omitempty"`
	Events           map[string]*Event           `json:"events,omitempty" yaml:"events,omitempty"`
}

type DataSource struct {
	Id           string            `json:"id" yaml:"id"`
	Name         string            `json:"name" yaml:"name"`
	Type         string            `json:"type" yaml:"type"`
	Path         string            `json:"path" yaml:"path"`
	Hash         string            `json:"hash,omitempty" yaml:"hash,omitempty"`
	ResourceName string            `json:"resourceName" yaml:"resourceName"`
	Description  string            `json:"description,omitempty" yaml:"description,omitempty"`
	Labels       map[string]string `json:"labels,omitempty" yaml:"labels,omitempty"`
}

type StoredProcedure struct {
	Metadata Metadata `json:"metadata" yaml:"metadata"`
	Control  Control  `json:"control" yaml:"control"`
	Features Features `json:"features" yaml:"features"`
	Links    Links    `json:"links" yaml:"links"`
}

type EventTrigger struct {
	Metadata Metadata `json:"metadata" yaml:"metadata"`
	Control  Control  `json:"control" yaml:"control"`
	Features Features `json:"features" yaml:"features"`
	Links    Links    `json:"links" yaml:"links"`
}

type Event struct {
	Metadata Metadata `json:"metadata" yaml:"metadata"`
	Control  Control  `json:"control" yaml:"control"`
	Features Features `json:"features" yaml:"features"`
}

type Metadata struct {
	Id          string            `json:"id" yaml:"id"`
	Name        string            `json:"name" yaml:"name"`
	Image       string            `json:"image,omitempty" yaml:"image,omitempty"`
	Build       Build             `json:"build,omitempty" yaml:"build,omitempty"`
	Hash        string            `json:"hash,omitempty" yaml:"hash,omitempty"`
	Prefix      string            `json:"prefix,omitempty" yaml:"prefix,omitempty"`
	Topic       string            `json:"topic,omitempty" yaml:"topic,omitempty"`
	Description string            `json:"description,omitempty" yaml:"description,omitempty"`
	Labels      map[string]string `json:"labels,omitempty" yaml:"labels,omitempty"`
	TriggerHash string            `json:"triggerHash,omitempty" yaml:"triggerHash,omitempty"`
}

type Build struct {
	Pull    string `json:"pull" yaml:"pull"`
	Workdir string `json:"workdir" yaml:"workdir"`
	Command string `json:"command" yaml:"command"`
}

type Control struct {
	DisableVirtualization bool   `json:"disableVirtualization" yaml:"disableVirtualization"`
	RunDetached           bool   `json:"runDetached" yaml:"runDetached"`
	RemoveOnStop          bool   `json:"removeOnStop" yaml:"removeOnStop"`
	Memory                string `json:"memory,omitempty" yaml:"memory,omitempty"`
	KernelArgs            string `json:"kernelArgs,omitempty" yaml:"kernelArgs,omitempty"`
}

type Features struct {
	Networks []string `json:"networks,omitempty" yaml:"networks,omitempty"`
	Ports    []string `json:"ports,omitempty" yaml:"ports,omitempty"`
	Volumes  []string `json:"volumes,omitempty" yaml:"volumes,omitempty"`
	Targets  []string `json:"targets,omitempty" yaml:"targets,omitempty"`
	EnvVars  []string `json:"envVars,omitempty" yaml:"envVars,omitempty"`
}

type Links struct {
	SoftLinks  []string `json:"softLinks,omitempty" yaml:"softLinks,omitempty"`
	HardLinks  []string `json:"hardLinks,omitempty" yaml:"hardLinks,omitempty"`
	EventLinks []string `json:"eventLinks,omitempty" yaml:"eventLinks,omitempty"`
}
