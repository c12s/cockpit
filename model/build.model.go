package model

// BuildSpec is the JSON form of supernova.ImageSpec. The gateway forwards it
// as the `spec` field of the BuildImage RPC.
type BuildSpec struct {
	Name          string            `json:"name"`
	Tag           string            `json:"tag,omitempty"`
	Containerfile string            `json:"containerfile,omitempty"`
	Registry      string            `json:"registry,omitempty"`
	BuildArgs     map[string]string `json:"buildArgs,omitempty"`
}

// BuildImageRequest is the JSON body of POST /apis/core/v1/builds. Context is
// the base64-encoded tar.gz of the build context; grpcurl decodes base64 into
// the proto bytes field.
type BuildImageRequest struct {
	Spec    BuildSpec `json:"spec"`
	Context string    `json:"context"`
}

// BuildQCOW2Request is the JSON body of POST /apis/core/v1/builds/qcow2.
type BuildQCOW2Request struct {
	Name       string `json:"name"`
	ImageRef   string `json:"imageRef"`
	TargetArch string `json:"targetArch,omitempty"`
	RootFS     string `json:"rootfs,omitempty"`
}

// BuildJob is the immediate reply to a build request.
type BuildJob struct {
	JobID string `json:"jobId"`
}

// BuildJobStatus is the JSON form of supernova.Job.
type BuildJobStatus struct {
	ID           string `json:"id"`
	Type         string `json:"type"`
	Status       string `json:"status"`
	ImageRef     string `json:"imageRef,omitempty"`
	Digest       string `json:"digest,omitempty"`
	ArtifactName string `json:"artifactName,omitempty"`
	ArtifactSum  string `json:"artifactSum,omitempty"`
	Error        string `json:"error,omitempty"`
	Log          string `json:"log,omitempty"`
}

// BuildJobsResponse is the JSON form of supernova.JobList.
type BuildJobsResponse struct {
	Jobs []BuildJobStatus `json:"jobs"`
}
