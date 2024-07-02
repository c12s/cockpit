package model

type ObjectScope struct {
	ID   string `json:"id" yaml:"id"`
	Kind string `json:"kind" yaml:"kind"`
}

type Permission struct {
	Condition struct {
		Expression string `json:"expression" yaml:"expression"`
	} `json:"condition" yaml:"condition"`
	Kind string `json:"kind" yaml:"kind"`
	Name string `json:"name" yaml:"name"`
}

type SubjectScope struct {
	ID   string `json:"id" yaml:"id"`
	Kind string `json:"kind" yaml:"kind"`
}

type PoliciesRequest struct {
	ObjectScope  ObjectScope  `json:"objectScope" yaml:"objectScope"`
	Permission   Permission   `json:"permission" yaml:"permission"`
	SubjectScope SubjectScope `json:"subjectScope" yaml:"subjectScope"`
}
