package model

// ImageListResponse is the JSON form of supernova.ImageList. It is the
// read-only view of the registry and the artifact server, returned by the
// gateway's ListImages route.
type ImageListResponse struct {
	Images    []string   `json:"images"`
	Artifacts []Artifact `json:"artifacts"`
}

// Artifact is one qcow2 on the artifact server, with its sha256.
type Artifact struct {
	Name     string `json:"name"`
	Checksum string `json:"checksum"`
}
