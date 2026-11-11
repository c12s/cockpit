package model

// ImageListResponse is the JSON form of supernova.ImageList. It is the
// read-only view of the registry and the artifact server, returned by the
// gateway's ListImages route.
type ImageListResponse struct {
	Images    []ImageEntry `json:"images"`
	Artifacts []Artifact   `json:"artifacts"`
}

// ImageEntry is one registry tag and the manifest digest it points to now,
// the value `image update --digest` takes.
type ImageEntry struct {
	Ref    string `json:"ref"`
	Digest string `json:"digest"`
}

// Artifact is one qcow2 on the artifact server, with its sha256.
type Artifact struct {
	Name     string `json:"name"`
	Checksum string `json:"checksum"`
}
