package model

import "time"

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

type HTTPRequestConfig struct {
	URL         string
	Method      string
	Headers     map[string]string
	RequestBody interface{}
	Response    interface{}
	Token       string
	Timeout     time.Duration
}
