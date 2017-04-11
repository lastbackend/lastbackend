package v1

import "time"

type Service struct {
	Name        string      `json:"name"`
	Description string      `json:"description"`
	Created     time.Time   `json:"created"`
	Updated     time.Time   `json:"updated"`
	Pods        *[]struct{} `json:"pods,omitempty"`
	Config      *Config     `json:"config,omitempty"`
}

type Config struct {
	Replicas int    `json:"scale,omitempty"`
	Memory   int    `json:"memory,omitempty"`
	Image    string `json:"image,omitempty"`
	Region   string `json:"region,omitempty"`
}

type ServiceList []Service
