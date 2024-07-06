package model

type SchemaDetails struct {
	Organization string `json:"organization" yaml:"organization"`
	SchemaName   string `json:"schemaName" yaml:"schemaName"`
	Version      string `json:"version,omitempty" yaml:"version,omitempty"`
	Namespace    string `json:"namespace,omitempty" yaml:"namespace,omitempty"`
}

type SchemaDetailsRequest struct {
	SchemaDetails SchemaDetails `json:"schema_details" yaml:"schema_details"`
}

type SchemaData struct {
	Schema       string `json:"schema" yaml:"schema"`
	CreationTime string `json:"creationTime" yaml:"creationTime"`
}

type SchemaVersion struct {
	SchemaDetails SchemaDetails `json:"schemaDetails" yaml:"schemaDetails"`
	SchemaData    SchemaData    `json:"schemaData" yaml:"schemaData"`
}

type SchemaVersionResponse struct {
	Message        string          `json:"message" yaml:"message"`
	SchemaVersions []SchemaVersion `json:"schemaVersions" yaml:"schemaVersions"`
}

type SchemaResponse struct {
	Message    string     `json:"message"`
	SchemaData SchemaData `json:"schemaData"`
}
