package model

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type RegistrationDetails struct {
	Email    string `json:"email"`
	Name     string `json:"name"`
	Org      string `json:"org"`
	Password string `json:"password"`
	Surname  string `json:"surname"`
	Username string `json:"username"`
}

type TokenResponse struct {
	Token string `json:"token"`
}

//
//type ConfigGroup struct {
//	Name    string            `yaml:"Name"`
//	OrgId   string            `yaml:"OrgId"`
//	Version int32             `yaml:"Version"`
//	Configs map[string]string `yaml:"Configs,omitempty"`
//}
//
//type Query struct {
//	Key      string `yaml:"Key"`
//	ShouldBe string `yaml:"ShouldBe"`
//	Value    string `yaml:"Value"`
//}
//
//type Policy struct {
//	SubjectScope Resource   `yaml:"SubjectScope"`
//	ObjectScope  Resource   `yaml:"ObjectScope"`
//	Permission   Permission `yaml:"Permission"`
//}
//
//type Resource struct {
//	Id   string `yaml:"Id"`
//	Kind string `yaml:"Kind"`
//}
//
//type Permission struct {
//	Name      string `yaml:"Name"`
//	Kind      string `yaml:"Kind"`
//	Condition string `yaml:"Condition"`
//}
